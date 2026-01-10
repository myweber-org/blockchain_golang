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
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{logFile: file}, nil
}

func (al *ActivityLogger) LogActivity(userID, action, resource string) error {
	activity := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}

	data, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = al.logFile.Write(data)
	return err
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
		userID   string
		action   string
		resource string
	}{
		{"user_001", "LOGIN", "auth_system"},
		{"user_002", "CREATE", "document_123"},
		{"user_001", "UPDATE", "profile_settings"},
		{"user_003", "DELETE", "comment_456"},
	}

	for _, act := range activities {
		err := logger.LogActivity(act.userID, act.action, act.resource)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		} else {
			fmt.Printf("Logged: %s performed %s on %s\n", act.userID, act.action, act.resource)
		}
	}
}