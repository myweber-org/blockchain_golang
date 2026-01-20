
package main

import (
	"log"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

type SessionStore interface {
	GetExpiredSessions() ([]Session, error)
	DeleteSession(id string) error
}

func cleanupExpiredSessions(store SessionStore) error {
	expiredSessions, err := store.GetExpiredSessions()
	if err != nil {
		return err
	}

	for _, session := range expiredSessions {
		err := store.DeleteSession(session.ID)
		if err != nil {
			log.Printf("Failed to delete session %s: %v", session.ID, err)
			continue
		}
		log.Printf("Deleted expired session: %s", session.ID)
	}

	return nil
}

func scheduleCleanup(store SessionStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := cleanupExpiredSessions(store)
		if err != nil {
			log.Printf("Session cleanup failed: %v", err)
		}
	}
}

func main() {
	// Implementation would provide concrete SessionStore
	var store SessionStore
	
	// Run cleanup every 24 hours
	scheduleCleanup(store, 24*time.Hour)
}