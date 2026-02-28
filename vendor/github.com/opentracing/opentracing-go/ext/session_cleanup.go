package main

import (
    "context"
    "log"
    "time"

    "yourproject/internal/db"
)

const cleanupInterval = 24 * time.Hour

func main() {
    ctx := context.Background()
    dbClient := db.NewClient()

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(ctx, dbClient); err != nil {
                log.Printf("Session cleanup failed: %v", err)
                continue
            }
            log.Println("Expired sessions cleaned up successfully")
        }
    }
}

func cleanupExpiredSessions(ctx context.Context, dbClient *db.Client) error {
    cutoff := time.Now().Add(-7 * 24 * time.Hour)
    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    
    result, err := dbClient.Exec(ctx, query, cutoff)
    if err != nil {
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    log.Printf("Deleted %d expired sessions", rowsAffected)
    return nil
}