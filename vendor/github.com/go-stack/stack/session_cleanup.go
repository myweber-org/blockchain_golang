
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
	"yourproject/internal/models"
)

func main() {
	ctx := context.Background()
	db := database.GetDB()

	for {
		now := time.Now()
		result := db.WithContext(ctx).Where("expires_at < ?", now).Delete(&models.Session{})
		
		if result.Error != nil {
			log.Printf("Error cleaning sessions: %v", result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Cleaned %d expired sessions", result.RowsAffected)
		}

		time.Sleep(24 * time.Hour)
	}
}