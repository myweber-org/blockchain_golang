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
}package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

type Authenticator struct {
	secretKey []byte
}

func NewAuthenticator(secretKey string) *Authenticator {
	return &Authenticator{
		secretKey: []byte(secretKey),
	}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
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

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.secretKey, nil
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

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "User ID not found in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}