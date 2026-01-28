
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
	ExpiresAt time.Time
	Data      map[string]interface{}
}

type Manager struct {
	sessions map[string]*Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
	}
}

func (m *Manager) CreateSession(userID int, duration time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
		Data:      make(map[string]interface{}),
	}

	m.sessions[token] = session
	return token, nil
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

	return session, nil
}

func (m *Manager) DeleteSession(token string) {
	delete(m.sessions, token)
}

func (m *Manager) SetSessionData(token string, key string, value interface{}) error {
	session, err := m.ValidateSession(token)
	if err != nil {
		return err
	}

	session.Data[key] = value
	return nil
}

func (m *Manager) GetSessionData(token string, key string) (interface{}, error) {
	session, err := m.ValidateSession(token)
	if err != nil {
		return nil, err
	}

	value, exists := session.Data[key]
	if !exists {
		return nil, errors.New("key not found in session data")
	}

	return value, nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}