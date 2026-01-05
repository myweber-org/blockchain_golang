
package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	sessionKeyPattern = "session:*"
	cleanupInterval   = 1 * time.Hour
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	ctx := context.Background()

	for range ticker.C {
		cleanupExpiredSessions(ctx, rdb)
	}
}

func cleanupExpiredSessions(ctx context.Context, rdb *redis.Client) {
	iter := rdb.Scan(ctx, 0, sessionKeyPattern, 0).Iterator()
	var keysToDelete []string

	for iter.Next(ctx) {
		key := iter.Val()
		ttl, err := rdb.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get TTL for key %s: %v", key, err)
			continue
		}
		if ttl < 0 {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Error during key scan: %v", err)
		return
	}

	if len(keysToDelete) > 0 {
		deleted, err := rdb.Del(ctx, keysToDelete...).Result()
		if err != nil {
			log.Printf("Failed to delete expired sessions: %v", err)
			return
		}
		log.Printf("Cleaned up %d expired sessions", deleted)
	}
}