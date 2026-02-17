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
    timeout  time.Duration
}

func NewManager(timeout time.Duration) *Manager {
    m := &Manager{
        sessions: make(map[string]*Session),
        timeout:  timeout,
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
        ExpiresAt: time.Now().Add(m.timeout),
    }
    m.sessions[session.ID] = session
    return session
}

func (m *Manager) Get(id string) (*Session, bool) {
    m.mu.RLock()
    session, exists := m.sessions[id]
    m.mu.RUnlock()

    if !exists {
        return nil, false
    }

    if time.Now().After(session.ExpiresAt) {
        m.Delete(id)
        return nil, false
    }
    return session, true
}

func (m *Manager) Delete(id string) {
    m.mu.Lock()
    delete(m.sessions, id)
    m.mu.Unlock()
}

func (m *Manager) cleanupLoop() {
    ticker := time.NewTicker(time.Hour)
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
    return "session_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
    }
    return string(b)
}