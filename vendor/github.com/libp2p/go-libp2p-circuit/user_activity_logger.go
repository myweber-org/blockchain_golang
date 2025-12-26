package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Activity struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) error {
	activity := Activity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now().UTC(),
		Details:   details,
	}

	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(activity); err != nil {
		return fmt.Errorf("failed to encode activity: %w", err)
	}

	return nil
}

func main() {
	activities := []struct {
		userID string
		action string
		details string
	}{
		{"user_001", "login", "successful authentication"},
		{"user_002", "purchase", "item_id: 456, amount: 29.99"},
		{"user_001", "logout", "session duration: 1h30m"},
	}

	for _, a := range activities {
		if err := logActivity(a.userID, a.action, a.details); err != nil {
			log.Printf("Failed to log activity: %v", err)
		} else {
			fmt.Printf("Logged %s action for user %s\n", a.action, a.userID)
		}
	}
}