package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	UserID    string
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

	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[token] = Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(m.ttl),
	}

	return token, nil
}

func (m *Manager) Validate(token string) (string, error) {
	m.mu.RLock()
	session, exists := m.sessions[token]
	m.mu.RUnlock()

	if !exists {
		return "", errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		m.mu.Lock()
		delete(m.sessions, token)
		m.mu.Unlock()
		return "", errors.New("session expired")
	}

	return session.UserID, nil
}

func (m *Manager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}package session

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    Data      map[string]interface{}
    ExpiresAt time.Time
    mu        sync.RWMutex
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

func (m *Manager) CreateSession() *Session {
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

func (m *Manager) GetSession(id string) *Session {
    m.mu.RLock()
    session, exists := m.sessions[id]
    m.mu.RUnlock()

    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
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

func (s *Session) Set(key string, value interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.Data[key] = value
}

func (s *Session) Get(key string) interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.Data[key]
}

func (s *Session) Extend(duration time.Duration) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.ExpiresAt = time.Now().Add(duration)
}

func generateID() string {
    return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
    }
    return string(b)
}