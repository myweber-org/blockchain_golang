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

func (a *Authenticator) ValidateToken(tokenString string) (bool, error) {
    if tokenString == "" {
        return false, nil
    }
    
    // In production, use proper JWT validation library
    // This is simplified for demonstration
    expectedPrefix := "valid_"
    return strings.HasPrefix(tokenString, expectedPrefix), nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }

        valid, err := a.ValidateToken(tokenParts[1])
        if err != nil {
            http.Error(w, "Token validation error", http.StatusInternalServerError)
            return
        }

        if !valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

func Authenticate(next http.Handler) http.Handler {
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
		userID, err := validateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func validateToken(token string) (string, error) {
	// Simplified token validation - in production use proper JWT library
	if token == "" || len(token) < 10 {
		return "", http.ErrAbortHandler
	}
	// Mock validation - return token as userID for demonstration
	return "user_" + token[:8], nil
}