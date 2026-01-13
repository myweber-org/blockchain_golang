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

type ActivityLogger struct {
	logFile *os.File
	encoder *json.Encoder
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{
		logFile: file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (al *ActivityLogger) LogActivity(userID, action, resource string) error {
	logEntry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}
	return al.encoder.Encode(logEntry)
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	activities := []struct {
		userID, action, resource string
	}{
		{"user123", "LOGIN", "auth_system"},
		{"user456", "CREATE", "document_789"},
		{"user123", "DOWNLOAD", "file_xyz.pdf"},
		{"user789", "UPDATE", "profile_settings"},
	}

	for _, act := range activities {
		if err := logger.LogActivity(act.userID, act.action, act.resource); err != nil {
			log.Printf("Failed to log activity: %v", err)
		} else {
			fmt.Printf("Logged: %s performed %s on %s\n", act.userID, act.action, act.resource)
		}
		time.Sleep(100 * time.Millisecond)
	}
}