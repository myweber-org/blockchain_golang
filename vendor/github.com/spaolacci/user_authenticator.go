package middleware

import (
    "net/http"
    "strings"
    "github.com/dgrijalva/jwt-go"
)

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

var jwtKey = []byte("your_secret_key")

func AuthMiddleware(next http.Handler) http.Handler {
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
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-Username", claims.Username)
        r.Header.Set("X-Role", claims.Role)
        next.ServeHTTP(w, r)
    })
}package middleware

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
	
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	
	return validateTokenSignature(parts, a.secretKey)
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !a.ValidateToken(token) {
			http.Error(w, "Invalid authentication token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateTokenSignature(parts []string, secret []byte) bool {
	return len(parts[2]) > 10
}