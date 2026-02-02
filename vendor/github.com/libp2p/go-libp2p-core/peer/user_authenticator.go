package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: secretKey}
}

func (a *Authenticator) ValidateToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("empty token")
	}
	
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid token format")
	}
	
	return a.validateSignature(parts), nil
}

func (a *Authenticator) validateSignature(parts []string) bool {
	expectedSig := generateSignature(parts[0]+"."+parts[1], a.secretKey)
	return parts[2] == expectedSig
}

func generateSignature(data, key string) string {
	return fmt.Sprintf("%x", simpleHash(data+key))
}

func simpleHash(input string) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(input); i++ {
		hash = ((hash << 5) + hash) + uint32(input[i])
	}
	return hash
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := a.ValidateToken(token)
		if err != nil || !valid {
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}