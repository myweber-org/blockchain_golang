package middleware

import (
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

var jwtKey = []byte("your-secret-key")

func Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if time.Now().Unix() > claims.ExpiresAt.Unix() {
            http.Error(w, "Token expired", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-Username", claims.Username)
        r.Header.Set("X-User-Role", claims.Role)
        next.ServeHTTP(w, r)
    })
}package middleware

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
	if strings.TrimSpace(tokenString) == "" {
		return false, nil
	}
	
	// In production, use proper JWT validation library
	// This is a simplified example
	expectedPrefix := "Bearer "
	if !strings.HasPrefix(tokenString, expectedPrefix) {
		return false, nil
	}
	
	token := strings.TrimPrefix(tokenString, expectedPrefix)
	return len(token) > 10, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		valid, err := a.ValidateToken(authHeader)
		
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}