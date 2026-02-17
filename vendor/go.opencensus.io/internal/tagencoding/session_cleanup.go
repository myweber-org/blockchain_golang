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
}