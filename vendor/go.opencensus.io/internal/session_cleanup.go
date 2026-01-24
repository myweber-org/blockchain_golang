package main

import (
    "context"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
}

func cleanupExpiredSessions() error {
    // Remove sessions older than 7 days
    cutoff := time.Now().Add(-7 * 24 * time.Hour).Unix()
    pattern := "session:*"

    iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
    for iter.Next(ctx) {
        key := iter.Val()
        created, err := rdb.HGet(ctx, key, "created_at").Int64()
        if err != nil {
            log.Printf("Failed to get creation time for %s: %v", key, err)
            continue
        }

        if created < cutoff {
            if err := rdb.Del(ctx, key).Err(); err != nil {
                log.Printf("Failed to delete expired session %s: %v", key, err)
            } else {
                log.Printf("Deleted expired session: %s", key)
            }
        }
    }

    if err := iter.Err(); err != nil {
        return err
    }

    return nil
}

func main() {
    initRedis()

    // Run cleanup daily at 2 AM
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            } else {
                log.Println("Session cleanup completed successfully")
            }
        }
    }
}