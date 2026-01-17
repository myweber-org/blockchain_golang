
package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/db"
    "github.com/yourproject/models"
)

const cleanupInterval = 24 * time.Hour

func main() {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    ctx := context.Background()
    database := db.GetConnection()

    for range ticker.C {
        err := cleanupExpiredSessions(ctx, database)
        if err != nil {
            log.Printf("Session cleanup failed: %v", err)
        } else {
            log.Println("Session cleanup completed successfully")
        }
    }
}

func cleanupExpiredSessions(ctx context.Context, db *db.Connection) error {
    cutoff := time.Now().Add(-models.SessionTTL)
    
    result, err := db.ExecContext(ctx,
        "DELETE FROM user_sessions WHERE last_activity < ?",
        cutoff,
    )
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Removed %d expired sessions", rowsAffected)
    return nil
}