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

	log.Println("Session cleanup service started")
	for range ticker.C {
		cleanupExpiredSessions(db)
	}
}

func cleanupExpiredSessions(db *database.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Failed to clean sessions: %v", err)
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned %d expired sessions", rows)
}