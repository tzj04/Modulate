package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"modulate/backend/internal/repositories"
	"modulate/backend/internal/models"
)

type UserHandler struct {
	Repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *UserHandler {
    return &UserHandler{Repo: repo}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
    Username string  `json:"username"`
    Password string  `json:"password"`
    Label    *string `json:"label"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Map request to Model
    user := models.User{
        Username: req.Username,
        Password: string(hashed),
        Label:    req.Label,
    }

    if err := h.Repo.Create(&user); err != nil {
        http.Error(w, "User already exists or DB error", http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

// Login validates user and returns a JWT
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := h.Repo.GetByUsername(req.Username)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Compare bcrypt hash
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Prepare Secret
    secret := []byte(os.Getenv("JWT_SECRET"))

    // Generate Short-Lived Access Token (1 hour)
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "iat":     time.Now().Unix(),
        "exp":     time.Now().Add(time.Hour).Unix(), 
    })

    // Create Long-Lived Refresh Token (7 days)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "iat":     time.Now().Unix(),
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
    })

    // Sign the tokens
    atString, err := accessToken.SignedString(secret)
    if err != nil {
        http.Error(w, "Error generating access token", http.StatusInternalServerError)
        return
    }
    
    rtString, err := refreshToken.SignedString(secret)
    if err != nil {
        http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
        return
    }

    // Send Refresh Token in an HttpOnly Cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    rtString,
        Expires:  time.Now().Add(time.Hour * 24 * 7),
        HttpOnly: true, 
        Secure:   false, // Set to true in production with HTTPS
        Path:     "/",
        SameSite: http.SameSiteLaxMode,
    })

    // Clean sensitive data
    user.Password = ""

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(AuthResponse{
        Token: atString,
        User:  *user,
    })
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
    // Get the refresh token from the cookie
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        http.Error(w, "Refresh token missing", http.StatusUnauthorized)
        return
    }

    // Parse and validate the refresh token
    secret := []byte(os.Getenv("JWT_SECRET"))
    token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
        return secret, nil
    })

    if err != nil || !token.Valid {
        http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
        return
    }

    // Extract user ID and issue a NEW access token
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Invalid claims", http.StatusUnauthorized)
        return
    }

    userID := int64(claims["user_id"].(float64))

    // New 15-minute Access Token
    newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "iat":     time.Now().Unix(),
        "exp":     time.Now().Add(time.Minute * 15).Unix(),
    })

    atString, err := newAccessToken.SignedString(secret)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Return the new access token
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": atString,
    })
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        Path:     "/",
    })
    w.WriteHeader(http.StatusOK)
}