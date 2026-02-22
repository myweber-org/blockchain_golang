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

func (a *Authenticator) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("empty token")
	}
	
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid token format")
	}
	
	return true, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if valid, err := a.ValidateToken(token); !valid {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}