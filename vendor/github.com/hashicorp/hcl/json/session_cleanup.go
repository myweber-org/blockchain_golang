package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	sessionKeyPattern = "session:*"
	batchSize         = 100
)

func cleanupExpiredSessions(ctx context.Context, client *redis.Client) error {
	var cursor uint64
	var keys []string
	var totalDeleted int64

	for {
		var err error
		keys, cursor, err = client.Scan(ctx, cursor, sessionKeyPattern, batchSize).Result()
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		if len(keys) > 0 {
			deleted, err := client.Del(ctx, keys...).Result()
			if err != nil {
				return fmt.Errorf("delete failed: %w", err)
			}
			totalDeleted += deleted
		}

		if cursor == 0 {
			break
		}
	}

	fmt.Printf("Cleaned up %d expired sessions\n", totalDeleted)
	return nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(ctx, rdb); err != nil {
				fmt.Printf("Cleanup error: %v\n", err)
			}
		}
	}
}