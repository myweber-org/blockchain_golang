
package main

import (
	"context"
	"log"
	"time"

	"your_project/internal/database"
)

const cleanupInterval = 24 * time.Hour
const sessionTTL = 7 * 24 * time.Hour

func main() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	ctx := context.Background()
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Session cleanup service started")

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx, db); err != nil {
				log.Printf("Cleanup failed: %v", err)
			} else {
				log.Println("Session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	cutoffTime := time.Now().Add(-sessionTTL)

	query := `DELETE FROM user_sessions WHERE last_activity < $1`
	result, err := db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Deleted %d expired sessions", rowsAffected)
	return nil
}