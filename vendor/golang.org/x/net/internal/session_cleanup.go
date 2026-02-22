
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

func cleanupExpiredSessions(ctx context.Context) error {
	db := database.GetDB()
	query := `DELETE FROM user_sessions WHERE expires_at < $1`
	result, err := db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
	return nil
}

func main() {
	ctx := context.Background()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}