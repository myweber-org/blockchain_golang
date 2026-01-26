package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	tokenValidator TokenValidator
	roleChecker    RoleChecker
}

type TokenValidator interface {
	Validate(token string) (UserClaims, error)
}

type RoleChecker interface {
	HasPermission(claims UserClaims, requiredRole string) bool
}

type UserClaims struct {
	UserID string
	Roles  []string
}

func NewAuthMiddleware(validator TokenValidator, checker RoleChecker) *AuthMiddleware {
	return &AuthMiddleware{
		tokenValidator: validator,
		roleChecker:    checker,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
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
		claims, err := m.tokenValidator.Validate(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("userClaims").(UserClaims)
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if !m.roleChecker.HasPermission(claims, requiredRole) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}