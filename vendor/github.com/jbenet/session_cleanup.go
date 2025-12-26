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

	for range ticker.C {
		if err := cleanupExpiredSessions(); err != nil {
			log.Printf("session cleanup failed: %v", err)
		} else {
			log.Println("session cleanup completed successfully")
		}
	}
}

func cleanupExpiredSessions() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.Conn.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	log.Printf("cleaned up %d expired sessions", rows)
	return nil
}