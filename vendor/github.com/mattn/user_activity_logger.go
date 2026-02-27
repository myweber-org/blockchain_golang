package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
}

func logActivity(userID, action, resource string) {
	activity := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}

	logEntry, err := json.Marshal(activity)
	if err != nil {
		log.Printf("Failed to marshal activity: %v", err)
		return
	}

	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(append(logEntry, '\n')); err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
}

func main() {
	logActivity("user123", "CREATE", "/api/document")
	logActivity("user456", "READ", "/api/report")
	logActivity("user789", "UPDATE", "/api/profile")

	fmt.Println("Activity logging completed")
}