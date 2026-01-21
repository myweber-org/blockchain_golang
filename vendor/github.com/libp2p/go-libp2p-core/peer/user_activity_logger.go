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
	Details   string    `json:"details,omitempty"`
}

type ActivityLogger struct {
	logs []ActivityLog
}

func NewActivityLogger() *ActivityLogger {
	return &ActivityLogger{
		logs: make([]ActivityLog, 0),
	}
}

func (al *ActivityLogger) LogActivity(userID, action, details string) {
	logEntry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
	al.logs = append(al.logs, logEntry)
	fmt.Printf("Logged: %s performed '%s' at %s\n", userID, action, logEntry.Timestamp.Format(time.RFC3339))
}

func (al *ActivityLogger) GetLogs() []ActivityLog {
	return al.logs
}

func (al *ActivityLogger) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(al.logs)
}

func (al *ActivityLogger) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&al.logs)
}

func main() {
	logger := NewActivityLogger()

	logger.LogActivity("user_123", "login", "Successful authentication")
	logger.LogActivity("user_456", "file_upload", "uploaded profile.jpg")
	logger.LogActivity("user_123", "logout", "Session terminated")

	if err := logger.SaveToFile("activity_log.json"); err != nil {
		log.Fatal("Failed to save logs:", err)
	}

	fmt.Println("Total activities logged:", len(logger.GetLogs()))
}