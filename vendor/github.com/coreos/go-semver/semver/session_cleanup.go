package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/database"
)

func main() {
    ctx := context.Background()
    db, err := database.NewConnection()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions(ctx, db)
        }
    }
}

func cleanupExpiredSessions(ctx context.Context, db *database.DB) {
    query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
    result, err := db.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Failed to clean up sessions: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
}