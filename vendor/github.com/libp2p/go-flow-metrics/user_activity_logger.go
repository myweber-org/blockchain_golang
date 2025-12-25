package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Activity struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
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

func (l *ActivityLogger) LogActivity(userID, action, details string) error {
	activity := Activity{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	data, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = l.logFile.Write(data)
	return err
}

func (l *ActivityLogger) Close() error {
	return l.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	activities := []struct {
		userID  string
		action  string
		details string
	}{
		{"user_001", "LOGIN", "User logged in from IP 192.168.1.100"},
		{"user_001", "VIEW_PAGE", "Accessed dashboard page"},
		{"user_002", "REGISTER", "New user registration completed"},
		{"user_001", "LOGOUT", "User session terminated"},
	}

	for _, act := range activities {
		err := logger.LogActivity(act.userID, act.action, act.details)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}

	fmt.Println("Activity logging completed")
}