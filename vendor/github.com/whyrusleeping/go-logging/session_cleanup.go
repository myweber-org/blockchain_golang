package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	cleanupInterval = 1 * time.Hour
	sessionTTL      = 7 * 24 * time.Hour
	deleteBatchSize = 1000
)

func cleanupExpiredSessions(db *pgxpool.Pool) error {
	ctx := context.Background()
	expiryTime := time.Now().Add(-sessionTTL)

	query := `
		DELETE FROM user_sessions 
		WHERE last_activity < $1 
		LIMIT $2
		RETURNING session_id
	`

	deletedIDs := make([]string, 0, deleteBatchSize)
	rows, err := db.Query(ctx, query, expiryTime, deleteBatchSize)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		deletedIDs = append(deletedIDs, id)
	}

	if len(deletedIDs) > 0 {
		log.Printf("Cleaned up %d expired sessions", len(deletedIDs))
	}
	return nil
}

func startSessionCleanup(db *pgxpool.Pool) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(db); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}

func main() {
	connStr := "postgresql://user:pass@localhost/dbname"
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	startSessionCleanup(db)
}