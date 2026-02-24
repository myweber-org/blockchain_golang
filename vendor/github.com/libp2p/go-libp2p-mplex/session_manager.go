package main

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "sync"
    "time"
)

type Session struct {
    UserID    string
    Token     string
    ExpiresAt time.Time
}

type SessionManager struct {
    sessions map[string]Session
    mu       sync.RWMutex
    duration time.Duration
}

func NewSessionManager(sessionDuration time.Duration) *SessionManager {
    return &SessionManager{
        sessions: make(map[string]Session),
        duration: sessionDuration,
    }
}

func generateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func (sm *SessionManager) CreateSession(userID string) (string, error) {
    token, err := generateToken()
    if err != nil {
        return "", err
    }

    session := Session{
        UserID:    userID,
        Token:     token,
        ExpiresAt: time.Now().Add(sm.duration),
    }

    sm.mu.Lock()
    sm.sessions[token] = session
    sm.mu.Unlock()

    return token, nil
}

func (sm *SessionManager) ValidateSession(token string) (string, error) {
    sm.mu.RLock()
    session, exists := sm.sessions[token]
    sm.mu.RUnlock()

    if !exists {
        return "", errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        sm.mu.Lock()
        delete(sm.sessions, token)
        sm.mu.Unlock()
        return "", errors.New("session expired")
    }

    return session.UserID, nil
}

func (sm *SessionManager) RemoveSession(token string) {
    sm.mu.Lock()
    delete(sm.sessions, token)
    sm.mu.Unlock()
}

func (sm *SessionManager) CleanupExpiredSessions() {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    now := time.Now()
    for token, session := range sm.sessions {
        if now.After(session.ExpiresAt) {
            delete(sm.sessions, token)
        }
    }
}