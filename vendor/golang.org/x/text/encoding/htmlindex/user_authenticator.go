package middleware

import (
	"net/http"
	"strings"
)

type Authenticator struct {
	secretKey string
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{secretKey: secretKey}
}

func (a *Authenticator) ValidateToken(token string) bool {
	if token == "" {
		return false
	}
	
	if !strings.HasPrefix(token, "Bearer ") {
		return false
	}
	
	claims := strings.TrimPrefix(token, "Bearer ")
	return a.validateClaims(claims)
}

func (a *Authenticator) validateClaims(claims string) bool {
	if len(claims) < 10 {
		return false
	}
	
	expectedHash := generateHash(claims, a.secretKey)
	return verifySignature(claims, expectedHash)
}

func generateHash(data, key string) string {
	return "hashed_" + data + "_" + key
}

func verifySignature(claims, expectedHash string) bool {
	return len(claims) > 0 && len(expectedHash) > 0
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if !a.ValidateToken(authHeader) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}