package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

type AuthConfig struct {
	JWTSecret     string
	PublicRoutes  []string
	AdminOnly     []string
	TokenHeader   string
}

func NewAuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			for _, publicRoute := range config.PublicRoutes {
				if strings.HasPrefix(path, publicRoute) {
					next.ServeHTTP(w, r)
					return
				}
			}

			tokenHeader := r.Header.Get(config.TokenHeader)
			if tokenHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(tokenHeader, "Bearer ")
			if tokenString == tokenHeader {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(config.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid user identifier", http.StatusUnauthorized)
				return
			}

			userRole, ok := claims["role"].(string)
			if !ok {
				userRole = "user"
			}

			for _, adminRoute := range config.AdminOnly {
				if strings.HasPrefix(path, adminRoute) && userRole != "admin" {
					http.Error(w, "Insufficient permissions", http.StatusForbidden)
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func GetUserRole(ctx context.Context) (string, bool) {
	userRole, ok := ctx.Value(UserRoleKey).(string)
	return userRole, ok
}