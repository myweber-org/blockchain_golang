package main

import (
	"context"
	"log"
	"time"
)

type SessionStore interface {
	DeleteExpired(ctx context.Context) error
}

type SessionCleaner struct {
	store     SessionStore
	interval  time.Duration
}

func NewSessionCleaner(store SessionStore, interval time.Duration) *SessionCleaner {
	return &SessionCleaner{
		store:    store,
		interval: interval,
	}
}

func (sc *SessionCleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(sc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Session cleaner stopped")
			return
		case <-ticker.C:
			if err := sc.store.DeleteExpired(ctx); err != nil {
				log.Printf("Failed to clean expired sessions: %v", err)
			} else {
				log.Println("Expired sessions cleaned successfully")
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := NewMemorySessionStore()
	cleaner := NewSessionCleaner(store, 1*time.Hour)

	go cleaner.Start(ctx)

	<-ctx.Done()
}