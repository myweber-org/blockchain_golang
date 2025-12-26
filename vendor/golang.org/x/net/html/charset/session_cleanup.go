package main

import (
    "log"
    "time"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
    dbConnection = "user=postgres dbname=appdb sslmode=disable"
    cleanupInterval = 24 * time.Hour
)

func cleanupExpiredSessions(db *sql.DB) error {
    query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
    result, err := db.Exec(query)
    if err != nil {
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func main() {
    db, err := sql.Open("postgres", dbConnection)
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
    defer db.Close()

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    log.Println("Session cleanup service started")
    
    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(db); err != nil {
                log.Printf("Cleanup error: %v", err)
            }
        }
    }
}