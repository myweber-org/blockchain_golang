package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

type Authenticator struct {
	secretKey []byte
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: []byte(secretKey)}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
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
		userID, err := a.validateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) validateToken(tokenString string) (string, error) {
	// Simplified token validation logic
	// In production, use proper JWT library like github.com/golang-jwt/jwt/v5
	if tokenString == "valid_test_token_123" {
		return "user_001", nil
	}
	return "", http.ErrNoCookie
}