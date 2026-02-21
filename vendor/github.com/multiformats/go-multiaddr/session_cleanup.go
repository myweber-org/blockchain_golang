
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
	deletedCount := 0

	for _, session := range sessions {
		if session.ExpiresAt.Before(now) {
			err := sc.store.DeleteSession(session.ID)
			if err != nil {
				log.Printf("Failed to delete session %s: %v", session.ID, err)
				continue
			}
			deletedCount++
		}
	}

	log.Printf("Cleaned up %d expired sessions", deletedCount)
	return nil
}

func (sc *SessionCleaner) StartDailyCleanup() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := sc.CleanExpiredSessions(); err != nil {
				log.Printf("Error cleaning sessions: %v", err)
			}
		}
	}
}

func main() {
	// Implementation would provide actual SessionStore
	var store SessionStore
	cleaner := NewSessionCleaner(store)
	
	go cleaner.StartDailyCleanup()
	
	// Keep main running
	select {}
}