package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type SessionCleaner struct {
    db        *sql.DB
    interval  time.Duration
    retention time.Duration
}

func NewSessionCleaner(db *sql.DB, interval, retention time.Duration) *SessionCleaner {
    return &SessionCleaner{
        db:        db,
        interval:  interval,
        retention: retention,
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
            sc.cleanExpiredSessions()
        }
    }
}

func (sc *SessionCleaner) cleanExpiredSessions() {
    cutoff := time.Now().Add(-sc.retention)
    query := `DELETE FROM user_sessions WHERE last_activity < $1`

    result, err := sc.db.Exec(query, cutoff)
    if err != nil {
        log.Printf("Failed to clean sessions: %v", err)
        return
    }

    rows, _ := result.RowsAffected()
    if rows > 0 {
        log.Printf("Cleaned %d expired sessions", rows)
    }
}package main

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

type SessionCleaner struct {
    db *sql.DB
}

func NewSessionCleaner(db *sql.DB) *SessionCleaner {
    return &SessionCleaner{db: db}
}

func (sc *SessionCleaner) Run(ctx context.Context) {
    ticker := time.NewTicker(cleanupInterval)
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

    cutoffTime := time.Now().Add(-sessionTTL)
    query := `DELETE FROM user_sessions WHERE last_activity < $1`

    result, err := sc.db.ExecContext(ctx, query, cutoffTime)
    if err != nil {
        log.Printf("Failed to clean expired sessions: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected > 0 {
        log.Printf("Cleaned %d expired sessions", rowsAffected)
    }
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    cleaner := NewSessionCleaner(db)
    ctx := context.Background()
    cleaner.Run(ctx)
}