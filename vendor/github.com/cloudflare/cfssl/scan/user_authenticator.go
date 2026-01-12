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
	return a.verifySignature(claims)
}

func (a *Authenticator) verifySignature(claims string) bool {
	expectedSignature := generateSignature(claims, a.secretKey)
	providedSignature := extractSignature(claims)
	
	return expectedSignature == providedSignature
}

func generateSignature(data, key string) string {
	return "sig_" + data[len(data)-8:] + "_" + key[:4]
}

func extractSignature(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
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