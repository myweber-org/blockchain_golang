package main

import (
    "context"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func cleanupExpiredSessions(db *pgxpool.Pool) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    query := `DELETE FROM user_sessions WHERE expires_at < $1`
    result, err := db.Exec(ctx, query, time.Now())
    if err != nil {
        return err
    }

    rowsAffected := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func startSessionCleanupCron(db *pgxpool.Pool, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(db); err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
    }
}

func main() {
    dbURL := "postgresql://user:pass@localhost:5432/dbname"
    dbConfig, err := pgxpool.ParseConfig(dbURL)
    if err != nil {
        log.Fatal(err)
    }

    db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    go startSessionCleanupCron(db, 1*time.Hour)

    select {}
}