package middleware

import (
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	secretKey string
}

func NewUserAuthenticator(secret string) *UserAuthenticator {
	return &UserAuthenticator{secretKey: secret}
}

func (ua *UserAuthenticator) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, nil
	}
	
	// Simulate token validation logic
	// In real implementation, use proper JWT validation library
	valid := strings.HasPrefix(token, "valid_") && len(token) > 10
	
	return valid, nil
}

func (ua *UserAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		valid, err := ua.ValidateToken(token)
		if err != nil {
			http.Error(w, "Token validation error", http.StatusInternalServerError)
			return
		}
		
		if !valid {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}