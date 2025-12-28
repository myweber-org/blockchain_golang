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

func (ua *UserAuthenticator) ValidateToken(token string) (string, error) {
	if token == "" {
		return "", http.ErrNoCookie
	}
	
	claims, err := parseJWT(token, ua.secretKey)
	if err != nil {
		return "", err
	}
	
	return claims.UserID, nil
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
		
		userID, err := ua.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		
		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	})
}

func parseJWT(token, secretKey string) (*TokenClaims, error) {
	// JWT parsing implementation would go here
	// This is a simplified placeholder
	return &TokenClaims{UserID: "sample-user-id"}, nil
}

type TokenClaims struct {
	UserID string
}