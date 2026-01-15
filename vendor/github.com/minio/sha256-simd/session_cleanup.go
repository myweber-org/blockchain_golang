
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	cleanupInterval := 24 * time.Hour

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	log.Println("Session cleanup service started")

	for {
		select {
		case <-ticker.C:
			err := cleanupExpiredSessions(ctx, db)
			if err != nil {
				log.Printf("Session cleanup failed: %v", err)
			} else {
				log.Println("Expired sessions cleaned successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	_, err := db.ExecContext(ctx, query)
	return err
}