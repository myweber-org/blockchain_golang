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

type SessionStore interface {
    GetExpiredSessions() ([]Session, error)
    DeleteSession(id string) error
}

func cleanupExpiredSessions(store SessionStore) error {
    expired, err := store.GetExpiredSessions()
    if err != nil {
        return err
    }

    for _, session := range expired {
        if err := store.DeleteSession(session.ID); err != nil {
            log.Printf("Failed to delete session %s: %v", session.ID, err)
        } else {
            log.Printf("Deleted expired session %s for user %s", session.ID, session.UserID)
        }
    }

    return nil
}

func scheduleCleanup(store SessionStore, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(store); err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
    }
}

func main() {
    // Implementation would provide actual SessionStore
    // store := NewDatabaseSessionStore()
    // scheduleCleanup(store, 24*time.Hour)
}