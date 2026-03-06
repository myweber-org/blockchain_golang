
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
			err := store.DeleteSession(session.ID)
			if err != nil {
				log.Printf("Failed to delete session %s: %v", session.ID, err)
			} else {
				log.Printf("Deleted expired session %s for user %s", session.ID, session.UserID)
			}
		}
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
	// In a real application, you would initialize your actual session store
	// For example: store := NewDatabaseSessionStore()
	// scheduleCleanup(store, 24*time.Hour)
	
	log.Println("Session cleanup scheduler started")
	select {} // Keep the program running
}