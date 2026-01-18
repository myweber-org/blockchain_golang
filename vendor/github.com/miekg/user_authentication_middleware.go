
package middleware

import (
	"net/http"
	"strings"
)

type JWTValidator interface {
	ValidateToken(tokenString string) (bool, error)
}

type AuthMiddleware struct {
	validator JWTValidator
}

func NewAuthMiddleware(validator JWTValidator) *AuthMiddleware {
	return &AuthMiddleware{validator: validator}
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
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
		valid, err := am.validator.ValidateToken(token)
		if err != nil || !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}