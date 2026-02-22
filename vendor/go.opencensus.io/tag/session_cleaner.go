
package main

import (
	"context"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Session struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	Data      []byte    `gorm:"type:jsonb"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func cleanupExpiredSessions(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&Session{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired sessions", result.RowsAffected)
	}

	return nil
}

func startCleanupJob(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := cleanupExpiredSessions(db); err != nil {
			log.Printf("Session cleanup failed: %v", err)
		}
	}
}

func main() {
	dsn := "host=localhost user=postgres password=secret dbname=app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&Session{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Starting session cleanup job...")
	startCleanupJob(db, 1*time.Hour)
}