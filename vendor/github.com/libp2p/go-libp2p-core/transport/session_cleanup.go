
package main

import (
	"context"
	"log"
	"time"
)

type SessionStore interface {
	DeleteExpiredSessions(ctx context.Context) error
}

type CleanupJob struct {
	store     SessionStore
	interval  time.Duration
}

func NewCleanupJob(store SessionStore, interval time.Duration) *CleanupJob {
	return &CleanupJob{
		store:    store,
		interval: interval,
	}
}

func (j *CleanupJob) Run(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Cleanup job stopped")
			return
		case <-ticker.C:
			if err := j.store.DeleteExpiredSessions(ctx); err != nil {
				log.Printf("Failed to delete expired sessions: %v", err)
			} else {
				log.Println("Expired sessions cleaned up successfully")
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := NewMemorySessionStore()
	job := NewCleanupJob(store, 24*time.Hour)

	go job.Run(ctx)

	<-ctx.Done()
}