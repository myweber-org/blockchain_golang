
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
)

const cleanupInterval = 24 * time.Hour

func main() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions()
		}
	}
}

func cleanupExpiredSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.Conn.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Failed to clean up sessions: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rowsAffected)
}