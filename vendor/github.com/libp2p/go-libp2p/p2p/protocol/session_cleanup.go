package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

func main() {
	ctx := context.Background()
	db := database.New()

	for {
		now := time.Now()
		cutoff := now.Add(-24 * time.Hour)

		result, err := db.ExecContext(ctx,
			"DELETE FROM user_sessions WHERE last_activity < ?",
			cutoff,
		)
		if err != nil {
			log.Printf("Failed to clean sessions: %v", err)
		} else {
			rows, _ := result.RowsAffected()
			log.Printf("Cleaned %d expired sessions", rows)
		}

		time.Sleep(24 * time.Hour)
	}
}