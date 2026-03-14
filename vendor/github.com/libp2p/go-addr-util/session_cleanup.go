package main

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

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions?sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(db); err != nil {
            log.Printf("Cleanup failed: %v", err)
        }
    }
}

func cleanupExpiredSessions(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cutoff := time.Now().Add(-sessionTTL)
    query := `DELETE FROM user_sessions WHERE last_activity < $1`

    result, err := db.ExecContext(ctx, query, cutoff)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rows)
    return nil
}