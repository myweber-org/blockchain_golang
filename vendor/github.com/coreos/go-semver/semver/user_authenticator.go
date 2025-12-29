package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secret string) *Authenticator {
	return &Authenticator{secretKey: secret}
}

func (a *Authenticator) ValidateToken(token string) bool {
	return token == fmt.Sprintf("valid-%s-token", a.secretKey)
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

		if !a.ValidateToken(parts[1]) {
			http.Error(w, "Invalid or expired token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}