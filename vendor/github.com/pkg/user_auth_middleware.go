package middleware

import (
    "net/http"
    "strings"
    "github.com/dgrijalva/jwt-go"
)

type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.StandardClaims
}

func JWTAuthMiddleware(secretKey string) func(http.Handler) http.Handler {
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

            tokenString := tokenParts[1]
            claims := &Claims{}

            token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(secretKey), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
                return
            }

            r.Header.Set("X-User-ID", claims.UserID)
            r.Header.Set("X-User-Role", claims.Role)

            next.ServeHTTP(w, r)
        })
    }
}