package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/db"
)

const cleanupInterval = 1 * time.Hour
const sessionTTL = 24 * time.Hour

func cleanupExpiredSessions(ctx context.Context) error {
    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    cutoffTime := time.Now().UTC().Add(-sessionTTL)

    result, err := db.Conn.ExecContext(ctx, query, cutoffTime)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func startSessionCleanup(ctx context.Context) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Session cleanup stopped")
            return
        case <-ticker.C:
            if err := cleanupExpiredSessions(ctx); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            }
        }
    }
}

func main() {
    ctx := context.Background()
    startSessionCleanup(ctx)
}