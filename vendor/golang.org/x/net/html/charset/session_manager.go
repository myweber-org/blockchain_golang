package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

var sessions = make(map[string]Session)

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateSession(userID int, duration time.Duration) (Session, error) {
	token, err := GenerateToken()
	if err != nil {
		return Session{}, err
	}

	now := time.Now()
	session := Session{
		ID:        token,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(duration),
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