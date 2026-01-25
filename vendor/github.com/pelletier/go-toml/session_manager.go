package session

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    UserID    int
    Data      map[string]interface{}
    ExpiresAt time.Time
}

type Manager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
    duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
    m := &Manager{
        sessions: make(map[string]*Session),
        duration: sessionDuration,
    }
    go m.cleanupLoop()
    return m
}

func (m *Manager) Create(userID int) *Session {
    m.mu.Lock()
    defer m.mu.Unlock()

    session := &Session{
        ID:        generateID(),
        UserID:    userID,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(m.duration),
    }
    m.sessions[session.ID] = session
    return session
}

func (m *Manager) Get(id string) *Session {
    m.mu.RLock()
    defer m.mu.RUnlock()

    session, exists := m.sessions[id]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
}

func (m *Manager) Refresh(id string) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    session, exists := m.sessions[id]
    if !exists {
        return false
    }
    session.ExpiresAt = time.Now().Add(m.duration)
    return true
}

func (m *Manager) cleanupLoop() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        m.cleanup()
    }
}

func (m *Manager) cleanup() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    for id, session := range m.sessions {
        if now.After(session.ExpiresAt) {
            delete(m.sessions, id)
        }
    }
}

func generateID() string {
    return "sess_" + time.Now().Format("20060102150405") + "_" + randomString(16)
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
    }
    return string(b)
}