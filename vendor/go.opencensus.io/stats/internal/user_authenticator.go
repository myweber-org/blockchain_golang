package auth

import (
	"context"
	"net/http"
	"strings"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func JWTMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			tokenStr := tokenParts[1]
			claims, err := validateToken(tokenStr, secretKey)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateToken(tokenStr string, secretKey string) (*Claims, error) {
	// Token validation logic would be implemented here
	// This is a placeholder for actual JWT validation
	return &Claims{
		UserID:   "sample_user_id",
		Username: "sample_user",
	}, nil
}