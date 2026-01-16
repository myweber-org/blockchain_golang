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

func GenerateSession(userID string, duration time.Duration) (Session, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return Session{}, err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(duration)

	session := Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	sessions[token] = session
	return session, nil
}

func ValidateSession(token string) (Session, error) {
	session, exists := sessions[token]
	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sessions, token)
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func InvalidateSession(token string) {
	delete(sessions, token)
}

func CleanupExpiredSessions() {
	now := time.Now()
	for token, session := range sessions {
		if now.After(session.ExpiresAt) {
			delete(sessions, token)
		}
	}
}