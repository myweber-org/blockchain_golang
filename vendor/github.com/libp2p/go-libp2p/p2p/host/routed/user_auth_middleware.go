package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	TokenService TokenValidator
}

type TokenValidator interface {
	ValidateToken(token string) (bool, error)
}

func NewAuthMiddleware(ts TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{TokenService: ts}
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		valid, err := am.TokenService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Token validation error", http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}