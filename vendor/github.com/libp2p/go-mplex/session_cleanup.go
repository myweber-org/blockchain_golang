
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
	GetAllSessions() ([]Session, error)
	DeleteSession(id string) error
}

type SessionCleaner struct {
	store SessionStore
}

func NewSessionCleaner(store SessionStore) *SessionCleaner {
	return &SessionCleaner{store: store}
}

func (sc *SessionCleaner) CleanExpiredSessions() error {
	sessions, err := sc.store.GetAllSessions()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, session := range sessions {
		if session.ExpiresAt.Before(now) {
			if err := sc.store.DeleteSession(session.ID); err != nil {
				log.Printf("Failed to delete session %s: %v", session.ID, err)
			} else {
				log.Printf("Deleted expired session %s for user %s", session.ID, session.UserID)
			}
		}
	}

	return nil
}

func (sc *SessionCleaner) StartCleanupJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := sc.CleanExpiredSessions(); err != nil {
			log.Printf("Session cleanup job failed: %v", err)
		}
	}
}

func main() {
	// This would be replaced with actual implementation
	var store SessionStore
	cleaner := NewSessionCleaner(store)

	// Run cleanup job daily
	go cleaner.StartCleanupJob(24 * time.Hour)

	// Keep main running
	select {}
}