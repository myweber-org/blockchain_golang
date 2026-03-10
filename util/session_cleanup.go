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

func cleanupExpiredSessions(db *sql.DB) error {
    ctx := context.Background()
    cutoff := time.Now().Add(-sessionTTL)

    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    result, err := db.ExecContext(ctx, query, cutoff)
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

func startSessionCleanup(db *sql.DB) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(db); err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
    }
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    startSessionCleanup(db)
}