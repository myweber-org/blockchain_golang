package main

import (
	"context"
	"log"
	"time"

	"github.com/yourproject/database"
)

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	cutoff := time.Now().Add(-24 * time.Hour)
	query := `DELETE FROM user_sessions WHERE last_activity < $1`
	result, err := db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
	return nil
}

func startSessionCleanupJob(db *database.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := cleanupExpiredSessions(ctx, db)
			cancel()
			if err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}

func main() {
	db := database.New()
	go startSessionCleanupJob(db)
	select {}
}