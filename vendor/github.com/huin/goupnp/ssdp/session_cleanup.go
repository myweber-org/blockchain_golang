package main

import (
    "log"
    "time"
    "database/sql"
    _ "github.com/lib/pq"
)

func cleanupExpiredSessions(db *sql.DB) error {
    query := `DELETE FROM user_sessions WHERE expires_at < $1`
    result, err := db.Exec(query, time.Now())
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func scheduleSessionCleanup(db *sql.DB) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(db); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            }
        }
    }
}

func main() {
    db, err := sql.Open("postgres", "host=localhost port=5432 user=app dbname=appdb sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    
    scheduleSessionCleanup(db)
}