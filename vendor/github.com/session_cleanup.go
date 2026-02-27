package main

import (
	"context"
	"log"
	"time"

	"github.com/yourproject/database"
)

func cleanupExpiredSessions(ctx context.Context, db *database.DB) error {
	cutoff := time.Now().Add(-24 * time.Hour)
	query := `DELETE FROM user_sessions WHERE last_activity < $1`
	result, err := db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
	return nil
}

func startSessionCleanupJob(db *database.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := cleanupExpiredSessions(ctx, db)
			cancel()
			if err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}

func main() {
	db := database.New()
	go startSessionCleanupJob(db)
	select {}
}package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const (
    cleanupInterval = 1 * time.Hour
    sessionTTL      = 24 * time.Hour
)

func cleanupExpiredSessions(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    cutoffTime := time.Now().Add(-sessionTTL)

    result, err := db.ExecContext(ctx, query, cutoffTime)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Failed to get rows affected: %v", err)
    } else {
        log.Printf("Cleaned up %d expired sessions", rowsAffected)
    }

    return nil
}

func startSessionCleanupJob(db *sql.DB) {
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
    db, err := sql.Open("postgres", "postgresql://user:pass@localhost/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    go startSessionCleanupJob(db)

    select {}
}