package middleware

import (
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secretKey: []byte(secret)}
}

func (am *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := parseJWTToken(tokenString, am.secretKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Expired() {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseJWTToken(tokenString string, secret []byte) (*TokenClaims, error) {
	// Token parsing implementation
	return validateTokenClaims(tokenString, secret)
}package middleware

import (
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	secretKey []byte
}

func NewUserAuthenticator(secret string) *UserAuthenticator {
	return &UserAuthenticator{secretKey: []byte(secret)}
}

func (ua *UserAuthenticator) ValidateToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return false, nil
	}

	return ua.verifySignature(parts), nil
}

func (ua *UserAuthenticator) verifySignature(parts []string) bool {
	expectedSig := generateSignature(parts[0]+"."+parts[1], ua.secretKey)
	return parts[2] == expectedSig
}

func generateSignature(data string, key []byte) string {
	hash := hmac.New(sha256.New, key)
	hash.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func (ua *UserAuthenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := ua.ValidateToken(token)
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
}