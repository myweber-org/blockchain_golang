package main

import (
    "log"
    "time"
)

type SessionStore struct {
    sessions map[string]time.Time
}

func (s *SessionStore) CleanupExpiredSessions(maxAge time.Duration) {
    cutoff := time.Now().Add(-maxAge)
    for sessionID, lastActivity := range s.sessions {
        if lastActivity.Before(cutoff) {
            delete(s.sessions, sessionID)
            log.Printf("Removed expired session: %s", sessionID)
        }
    }
}

func main() {
    store := &SessionStore{
        sessions: make(map[string]time.Time),
    }

    // Simulate some sessions
    store.sessions["abc123"] = time.Now().Add(-2 * time.Hour)
    store.sessions["def456"] = time.Now().Add(-30 * time.Minute)
    store.sessions["ghi789"] = time.Now().Add(-25 * time.Hour)

    log.Printf("Before cleanup: %d sessions", len(store.sessions))
    store.CleanupExpiredSessions(24 * time.Hour)
    log.Printf("After cleanup: %d sessions", len(store.sessions))
}