package middleware

import (
	"context"
	"net/http"
	"strings"
)

// Context key for user ID
type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware validates JWT token (placeholder) and injects userID into context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		token := tokenParts[1]

		// TODO: Replace with real JWT validation
		// For now, fake userID = 1 if token is "test"
		var userID int64
		if token == "test" {
			userID = 1
		} else {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
