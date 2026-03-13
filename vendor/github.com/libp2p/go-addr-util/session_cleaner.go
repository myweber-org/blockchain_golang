
package main

import (
	"log"
	"time"
)

type SessionStore struct {
	sessions map[string]time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]time.Time),
	}
}

func (s *SessionStore) CleanExpiredSessions(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)
	for sessionID, createdAt := range s.sessions {
		if createdAt.Before(cutoff) {
			delete(s.sessions, sessionID)
			log.Printf("Removed expired session: %s", sessionID)
		}
	}
}

func startSessionCleanup(store *SessionStore, interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		store.CleanExpiredSessions(maxAge)
	}
}

func main() {
	sessionStore := NewSessionStore()
	
	go startSessionCleanup(sessionStore, 24*time.Hour, 7*24*time.Hour)
	
	select {}
}