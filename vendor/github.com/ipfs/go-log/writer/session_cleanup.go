
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

func cleanupExpiredSessions(store SessionStore) error {
	sessions, err := store.GetAllSessions()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, session := range sessions {
		if session.ExpiresAt.Before(now) {
			if err := store.DeleteSession(session.ID); err != nil {
				log.Printf("Failed to delete session %s: %v", session.ID, err)
			} else {
				log.Printf("Deleted expired session %s for user %s", session.ID, session.UserID)
			}
		}
	}
	return nil
}

func scheduleSessionCleanup(store SessionStore) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(store); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}

func main() {
	// Implementation would initialize SessionStore
	// and call scheduleSessionCleanup()
}