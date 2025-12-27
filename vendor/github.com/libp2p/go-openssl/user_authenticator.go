package middleware

import (
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	secretKey string
}

func NewUserAuthenticator(secretKey string) *UserAuthenticator {
	return &UserAuthenticator{secretKey: secretKey}
}

func (ua *UserAuthenticator) ValidateToken(token string) (bool, string) {
	if token == "" {
		return false, ""
	}

	claims, err := parseJWTToken(token, ua.secretKey)
	if err != nil {
		return false, ""
	}

	return true, claims.UserID
}

func (ua *UserAuthenticator) Middleware(next http.Handler) http.Handler {
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

		token := parts[1]
		valid, userID := ua.ValidateToken(token)
		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	})
}