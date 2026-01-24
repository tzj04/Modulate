package middleware

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

// Context key for user ID
type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware validates JWT token and injects userID into context
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip auth for preflight requests
        if r.Method == http.MethodOptions {
            next.ServeHTTP(w, r)
            return
        }

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "missing Authorization header", http.StatusUnauthorized)
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
            return
        }

        tokenString := tokenParts[1]
        secret := []byte(os.Getenv("JWT_SECRET"))

        // Parse and validate the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Verify the signing method is HMAC
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return secret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "session expired or invalid token", http.StatusUnauthorized)
            return
        }

        // Extract claims and set userID in context
        var userID int64
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // user_id comes from the payload of the JWT created during login
            if uid, ok := claims["user_id"].(float64); ok {
                userID = int64(uid)
            } else {
                http.Error(w, "invalid token payload: user_id missing", http.StatusUnauthorized)
                return
            }
        } else {
            http.Error(w, "invalid token claims", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}