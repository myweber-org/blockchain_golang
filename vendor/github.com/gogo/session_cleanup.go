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
    batchSize       = 1000
)

func main() {
    dbPool, err := pgxpool.New(context.Background(), "postgresql://user:pass@localhost/db")
    if err != nil {
        log.Fatal("Failed to create connection pool:", err)
    }
    defer dbPool.Close()

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(dbPool); err != nil {
            log.Printf("Cleanup failed: %v", err)
        } else {
            log.Printf("Cleanup completed successfully")
        }
    }
}

func cleanupExpiredSessions(db *pgxpool.Pool) error {
    ctx := context.Background()
    cutoff := time.Now().Add(-sessionTTL)

    for {
        result, err := db.Exec(ctx,
            `DELETE FROM user_sessions 
             WHERE last_activity < $1 
             AND session_id IN (
                 SELECT session_id FROM user_sessions 
                 WHERE last_activity < $1 
                 LIMIT $2
             )`,
            cutoff, batchSize,
        )
        if err != nil {
            return err
        }

        rowsAffected := result.RowsAffected()
        if rowsAffected == 0 {
            break
        }
        log.Printf("Deleted %d expired sessions", rowsAffected)
    }
    return nil
}