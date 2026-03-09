package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
)

func main() {
	ctx := context.Background()
	dbClient, err := db.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := cleanupExpiredSessions(ctx, dbClient)
			if err != nil {
				log.Printf("Session cleanup failed: %v", err)
			} else {
				log.Println("Session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, dbClient *db.Client) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	_, err := dbClient.ExecContext(ctx, query)
	return err
}package main

import (
    "context"
    "log"
    "time"
)

type SessionStore interface {
    DeleteExpired(ctx context.Context) error
}

type CleanupService struct {
    store SessionStore
}

func NewCleanupService(store SessionStore) *CleanupService {
    return &CleanupService{store: store}
}

func (s *CleanupService) RunDailyCleanup(ctx context.Context) error {
    log.Println("Starting daily session cleanup")
    
    if err := s.store.DeleteExpired(ctx); err != nil {
        log.Printf("Cleanup failed: %v", err)
        return err
    }
    
    log.Println("Session cleanup completed successfully")
    return nil
}

func (s *CleanupService) StartScheduler(ctx context.Context) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            log.Println("Cleanup scheduler stopped")
            return
        case <-ticker.C:
            if err := s.RunDailyCleanup(ctx); err != nil {
                log.Printf("Scheduled cleanup error: %v", err)
            }
        }
    }
}

func main() {
    ctx := context.Background()
    
    store := &InMemorySessionStore{}
    service := NewCleanupService(store)
    
    go service.StartScheduler(ctx)
    
    <-ctx.Done()
}

type InMemorySessionStore struct{}

func (s *InMemorySessionStore) DeleteExpired(ctx context.Context) error {
    return nil
}