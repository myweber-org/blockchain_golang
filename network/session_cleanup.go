
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

func main() {
    db, err := sql.Open("postgres", "host=localhost port=5432 user=app dbname=appdb password=secret sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    for {
        if err := cleanupExpiredSessions(db); err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
        
        time.Sleep(24 * time.Hour)
    }
}