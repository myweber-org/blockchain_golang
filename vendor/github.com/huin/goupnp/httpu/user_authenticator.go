
package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("your-secret-key-here")

type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateToken(userID, username string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auth-service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
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
    
    // Token validation logic would go here
    // For this example, we'll simulate validation
    isValid := strings.HasPrefix(tokenString, "valid_")
    return isValid, nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        tokenParts := strings.Split(authHeader, "Bearer ")
        if len(tokenParts) != 2 {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }
        
        token := tokenParts[1]
        valid, err := a.ValidateToken(token)
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