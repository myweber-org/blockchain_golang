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

type SessionCleaner struct {
    db *pgxpool.Pool
}

func NewSessionCleaner(db *pgxpool.Pool) *SessionCleaner {
    return &SessionCleaner{db: db}
}

func (sc *SessionCleaner) Run(ctx context.Context) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            sc.cleanupExpiredSessions(ctx)
        }
    }
}

func (sc *SessionCleaner) cleanupExpiredSessions(ctx context.Context) {
    cutoffTime := time.Now().Add(-sessionTTL)
    deletedCount := 0

    for {
        result, err := sc.db.Exec(ctx,
            `DELETE FROM user_sessions 
             WHERE last_activity < $1 
             AND id IN (
                 SELECT id FROM user_sessions 
                 WHERE last_activity < $1 
                 LIMIT $2
             )`,
            cutoffTime, deleteBatchSize)

        if err != nil {
            log.Printf("Error deleting expired sessions: %v", err)
            return
        }

        rowsAffected := result.RowsAffected()
        deletedCount += int(rowsAffected)

        if rowsAffected < deleteBatchSize {
            break
        }

        time.Sleep(100 * time.Millisecond)
    }

    if deletedCount > 0 {
        log.Printf("Cleaned up %d expired user sessions", deletedCount)
    }
}

func main() {
    ctx := context.Background()
    
    dbConfig, err := pgxpool.ParseConfig("postgresql://user:pass@localhost/db")
    if err != nil {
        log.Fatal(err)
    }

    dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer dbPool.Close()

    cleaner := NewSessionCleaner(dbPool)
    cleaner.Run(ctx)
}