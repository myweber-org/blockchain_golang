package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		userID, err := validateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(tokenString string) (string, error) {
	// In a real implementation, this would parse and validate JWT
	// For this example, we'll do a simple mock validation
	if tokenString == "" || len(tokenString) < 10 {
		return "", http.ErrNoCookie
	}
	
	// Mock extraction of user ID from token
	// In reality, this would decode JWT and verify signature
	userID := "user_" + tokenString[:8]
	return userID, nil
}