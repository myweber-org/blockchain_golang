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
	Data      map[string]interface{}
}

type Manager struct {
	sessions map[string]*Session
	duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		duration: sessionDuration,
	}
}

func (m *Manager) CreateSession(userID int, initialData map[string]interface{}) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.duration),
		Data:      initialData,
	}

	m.sessions[token] = session
	return session, nil
}

func (m *Manager) ValidateSession(token string) (*Session, error) {
	session, exists := m.sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(m.sessions, token)
		return nil, errors.New("session expired")
	}

	session.ExpiresAt = time.Now().Add(m.duration)
	return session, nil
}

func (m *Manager) InvalidateSession(token string) {
	delete(m.sessions, token)
}

func (m *Manager) CleanupExpired() {
	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}