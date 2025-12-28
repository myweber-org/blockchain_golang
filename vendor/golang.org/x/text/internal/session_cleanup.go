package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type SessionCleaner struct {
    db        *sql.DB
    batchSize int
    interval  time.Duration
}

func NewSessionCleaner(db *sql.DB, batchSize int, interval time.Duration) *SessionCleaner {
    return &SessionCleaner{
        db:        db,
        batchSize: batchSize,
        interval:  interval,
    }
}

func (sc *SessionCleaner) Start(ctx context.Context) {
    ticker := time.NewTicker(sc.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Session cleaner stopped")
            return
        case <-ticker.C:
            sc.cleanupExpiredSessions()
        }
    }
}

func (sc *SessionCleaner) cleanupExpiredSessions() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    query := `DELETE FROM user_sessions WHERE expires_at < NOW() LIMIT $1`
    result, err := sc.db.ExecContext(ctx, query, sc.batchSize)
    if err != nil {
        log.Printf("Failed to clean sessions: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected > 0 {
        log.Printf("Cleaned %d expired sessions", rowsAffected)
    }
}