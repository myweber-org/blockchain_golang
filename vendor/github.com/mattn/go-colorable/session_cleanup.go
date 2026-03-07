package main

import (
    "log"
    "time"
)

type Session struct {
    ID        string
    UserID    string
    ExpiresAt time.Time
}

type SessionStore struct {
    sessions map[string]Session
}

func NewSessionStore() *SessionStore {
    return &SessionStore{
        sessions: make(map[string]Session),
    }
}

func (s *SessionStore) CleanExpiredSessions() {
    now := time.Now()
    expiredCount := 0
    
    for id, session := range s.sessions {
        if session.ExpiresAt.Before(now) {
            delete(s.sessions, id)
            expiredCount++
        }
    }
    
    if expiredCount > 0 {
        log.Printf("Cleaned %d expired sessions", expiredCount)
    }
}

func startSessionCleanupJob(store *SessionStore, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            store.CleanExpiredSessions()
        }
    }
}

func main() {
    sessionStore := NewSessionStore()
    
    // Add some test sessions
    sessionStore.sessions["test1"] = Session{
        ID:        "test1",
        UserID:    "user123",
        ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
    }
    
    sessionStore.sessions["test2"] = Session{
        ID:        "test2",
        UserID:    "user456",
        ExpiresAt: time.Now().Add(1 * time.Hour), // Still valid
    }
    
    // Start cleanup job running every 5 minutes
    go startSessionCleanupJob(sessionStore, 5*time.Minute)
    
    // Keep main goroutine alive
    select {}
}