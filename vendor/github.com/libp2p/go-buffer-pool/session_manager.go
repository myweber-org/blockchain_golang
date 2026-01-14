package session

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    Data      map[string]interface{}
    ExpiresAt time.Time
}

type Manager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
    ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
    m := &Manager{
        sessions: make(map[string]*Session),
        ttl:      ttl,
    }
    go m.cleanupLoop()
    return m
}

func (m *Manager) Create() *Session {
    m.mu.Lock()
    defer m.mu.Unlock()

    id := generateID()
    session := &Session{
        ID:        id,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(m.ttl),
    }
    m.sessions[id] = session
    return session
}

func (m *Manager) Get(id string) (*Session, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    session, exists := m.sessions[id]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil, false
    }
    return session, true
}

func (m *Manager) cleanupLoop() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        m.mu.Lock()
        now := time.Now()
        for id, session := range m.sessions {
            if now.After(session.ExpiresAt) {
                delete(m.sessions, id)
            }
        }
        m.mu.Unlock()
    }
}

func generateID() string {
    return time.Now().UnixNano()
}package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Manager struct {
	sessions map[string]Session
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]Session),
		ttl:      ttl,
	}
}

func (m *Manager) Create(userID string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.ttl),
	}

	m.mu.Lock()
	m.sessions[token] = session
	m.mu.Unlock()

	return token, nil
}

func (m *Manager) Validate(token string) (Session, error) {
	m.mu.RLock()
	session, exists := m.sessions[token]
	m.mu.RUnlock()

	if !exists {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		m.mu.Lock()
		delete(m.sessions, token)
		m.mu.Unlock()
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func (m *Manager) Invalidate(token string) {
	m.mu.Lock()
	delete(m.sessions, token)
	m.mu.Unlock()
}

func (m *Manager) Cleanup() {
	m.mu.Lock()
	for token, session := range m.sessions {
		if time.Now().After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
	m.mu.Unlock()
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}