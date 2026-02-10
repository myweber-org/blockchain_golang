
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

	ctx := context.Background()
	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx); err != nil {
				log.Printf("session cleanup failed: %v", err)
			} else {
				log.Println("session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.Conn.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	log.Printf("deleted %d expired sessions", rows)
	return nil
}
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

const cleanupInterval = 24 * time.Hour

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx, db); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			} else {
				log.Println("Session cleanup completed successfully")
			}
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Cleaned up %d expired sessions", rowsAffected)
	return nil
}