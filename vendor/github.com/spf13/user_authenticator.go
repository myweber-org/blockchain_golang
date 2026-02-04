package auth

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

var jwtKey = []byte("your_secret_key_here")

func GenerateToken(username, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
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

        if time.Until(claims.ExpiresAt.Time) < 0 {
            http.Error(w, "Token expired", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-Username", claims.Username)
        r.Header.Set("X-Role", claims.Role)
        next.ServeHTTP(w, r)
    }
}package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

var jwtKey = []byte("your_secret_key_here")

func GenerateToken(userID, email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("your_secret_key_here")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if time.Until(claims.ExpiresAt.Time) < 30*time.Second {
			http.Error(w, "Token expiring soon", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Username", claims.Username)

		next.ServeHTTP(w, r)
	})
}

func GenerateToken(userID, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth_service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       string
	Username string
	Email    string
	Role     string
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	jwtSecret     []byte
	tokenDuration time.Duration
}

func NewAuthService(secret string, duration time.Duration) (*AuthService, error) {
	if secret == "" {
		randomSecret, err := generateRandomSecret(32)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random secret: %w", err)
		}
		secret = randomSecret
	}

	return &AuthService{
		jwtSecret:     []byte(secret),
		tokenDuration: duration,
	}, nil
}

func generateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (a *AuthService) GenerateToken(user *User) (string, error) {
	expirationTime := time.Now().Add(a.tokenDuration)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (a *AuthService) RefreshToken(tokenString string) (string, error) {
	claims, err := a.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("cannot refresh invalid token: %w", err)
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return "", errors.New("token is not eligible for refresh yet")
	}

	newClaims := &Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   claims.UserID,
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err = newToken.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refreshed token: %w", err)
	}

	return tokenString, nil
}

func main() {
	authService, err := NewAuthService("", 15*time.Minute)
	if err != nil {
		fmt.Printf("Failed to create auth service: %v\n", err)
		return
	}

	user := &User{
		ID:       "12345",
		Username: "john_doe",
		Email:    "john@example.com",
		Role:     "user",
	}

	token, err := authService.GenerateToken(user)
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		return
	}

	fmt.Printf("Generated token: %s\n", token)

	claims, err := authService.ValidateToken(token)
	if err != nil {
		fmt.Printf("Token validation failed: %v\n", err)
		return
	}

	fmt.Printf("Token validated successfully. User: %s, Role: %s\n", claims.Username, claims.Role)

	refreshedToken, err := authService.RefreshToken(token)
	if err != nil {
		fmt.Printf("Token refresh failed: %v\n", err)
		return
	}

	fmt.Printf("Refreshed token: %s\n", refreshedToken)
}package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	secretKey       = []byte("your-secret-key-change-in-production")
)

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
			return nil, ErrInvalidToken
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Bearer token required", http.StatusUnauthorized)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-User-Role", claims.Role)
			next.ServeHTTP(w, r)
		})
	}
}package middleware

import (
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
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
            claims := &Claims{}

            token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(secretKey), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            r.Header.Set("X-User-ID", claims.UserID)
            r.Header.Set("X-User-Role", claims.Role)
            next.ServeHTTP(w, r)
        })
    }
}package middleware

import (
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

var jwtKey = []byte("your_secret_key_here")

func AuthMiddleware(next http.Handler) http.Handler {
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
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        if time.Until(claims.ExpiresAt.Time) < 30*time.Second {
            http.Error(w, "Token expiring soon", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-Username", claims.Username)
        r.Header.Set("X-UserRole", claims.Role)

        next.ServeHTTP(w, r)
    })
}

func GenerateToken(username, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auth-service",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}package auth

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const userIDKey contextKey = "userID"

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

        tokenString := parts[1]
        userID, err := a.validateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (a *Authenticator) validateToken(tokenString string) (string, error) {
    // Token validation logic would go here
    // For this example, we'll just check if it's a non-empty string
    if tokenString == "" {
        return "", http.ErrNoCookie
    }
    
    // In real implementation, this would parse and verify JWT
    // using the secretKey and return the user ID from claims
    return "user123", nil
}

func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}