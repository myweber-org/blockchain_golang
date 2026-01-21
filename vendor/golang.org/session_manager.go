package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

var sessions = make(map[string]Session)

func GenerateSession(userID int) (string, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	tokenStr := hex.EncodeToString(token)
	expiration := time.Now().Add(24 * time.Hour)

	sessions[tokenStr] = Session{
		ID:        tokenStr,
		UserID:    userID,
		ExpiresAt: expiration,
	}

	return tokenStr, nil
}

func ValidateSession(token string) (int, error) {
	session, exists := sessions[token]
	if !exists {
		return 0, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(sessions, token)
		return 0, errors.New("session expired")
	}

	return session.UserID, nil
}

func InvalidateSession(token string) {
	delete(sessions, token)
}