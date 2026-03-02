
package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/database"
)

const cleanupInterval = 24 * time.Hour
const sessionTTL = 7 * 24 * time.Hour

func cleanupExpiredSessions(ctx context.Context) error {
    db := database.GetConnection()
    cutoff := time.Now().Add(-sessionTTL)

    result, err := db.ExecContext(ctx,
        "DELETE FROM user_sessions WHERE last_activity < ?",
        cutoff)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func startSessionCleanupScheduler() {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            ctx := context.Background()
            if err := cleanupExpiredSessions(ctx); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            }
        }
    }
}

func main() {
    log.Println("Starting session cleanup scheduler...")
    startSessionCleanupScheduler()
}