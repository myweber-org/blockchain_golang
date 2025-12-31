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
    rows, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rows)
    return nil
}

func scheduleCleanup(db *sql.DB) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(db); err != nil {
                log.Printf("Cleanup failed: %v", err)
            }
        }
    }
}

func main() {
    db, err := sql.Open("postgres", "host=localhost user=app dbname=appdb sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    scheduleCleanup(db)
}