package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

type Session struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

var sessions = make(map[string]Session)

func GenerateSession(userID string) (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(24 * time.Hour)

	session := Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	sessions[token] = session
	return token, nil
}

func ValidateSession(token string) (string, error) {
	session, exists := sessions[token]
	if !exists {
		return "", errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sessions, token)
		return "", errors.New("session expired")
	}

	return session.UserID, nil
}

func InvalidateSession(token string) {
	delete(sessions, token)
}