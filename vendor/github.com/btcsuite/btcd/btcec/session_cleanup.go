
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

func cleanupExpiredSessions() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < $1`
	result, err := db.Pool.Exec(ctx, query, time.Now())
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rowsAffected)
	return nil
}

func main() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := cleanupExpiredSessions(); err != nil {
			log.Printf("Session cleanup failed: %v", err)
		}
	}
}