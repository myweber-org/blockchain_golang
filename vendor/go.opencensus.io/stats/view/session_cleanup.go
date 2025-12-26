
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
	"yourproject/internal/models"
)

func cleanupExpiredSessions() error {
	db := database.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	threshold := time.Now().Add(-24 * time.Hour)
	result := db.WithContext(ctx).Where("last_activity < ?", threshold).Delete(&models.Session{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired sessions older than %v", result.RowsAffected, threshold)
	}
	return nil
}

func main() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredSessions(); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}