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
	query := `DELETE FROM user_sessions WHERE expires_at < $1`
	result, err := db.ExecContext(ctx, query, time.Now())
	if err != nil {
		log.Printf("Failed to clean up sessions: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rowsAffected)
}