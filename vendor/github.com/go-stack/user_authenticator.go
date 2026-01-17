package middleware

import (
    "net/http"
    "strings"
)

type Authenticator struct {
    secretKey []byte
}

func NewAuthenticator(secret string) *Authenticator {
    return &Authenticator{secretKey: []byte(secret)}
}

func (a *Authenticator) ValidateToken(token string) bool {
    if token == "" {
        return false
    }
    
    if !strings.HasPrefix(token, "Bearer ") {
        return false
    }
    
    token = strings.TrimPrefix(token, "Bearer ")
    
    return len(token) > 10
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