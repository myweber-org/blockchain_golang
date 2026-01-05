package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]
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
	// Token validation logic would be implemented here
	// For example, using github.com/golang-jwt/jwt
	// This is a placeholder implementation
	if tokenString == "" {
		return "", http.ErrNoCookie
	}
	// In real implementation, parse and verify JWT
	return "user123", nil
}