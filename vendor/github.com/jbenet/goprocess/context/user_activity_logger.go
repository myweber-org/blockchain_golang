package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityType string

const (
	Login    ActivityType = "LOGIN"
	Logout   ActivityType = "LOGOUT"
	Purchase ActivityType = "PURCHASE"
	View     ActivityType = "VIEW"
)

type UserActivity struct {
	UserID    string       `json:"user_id"`
	Action    ActivityType `json:"action"`
	Timestamp time.Time    `json:"timestamp"`
	Details   string       `json:"details,omitempty"`
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

func (l *ActivityLogger) LogActivity(activity UserActivity) error {
	activity.Timestamp = time.Now().UTC()
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
	logger, err := NewActivityLogger("user_activities.jsonl")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	activities := []UserActivity{
		{UserID: "user123", Action: Login, Details: "Successful login from Chrome"},
		{UserID: "user123", Action: View, Details: "Viewed product page: laptop"},
		{UserID: "user123", Action: Purchase, Details: "Purchased item: laptop"},
		{UserID: "user456", Action: Login, Details: "Mobile app login"},
	}

	for _, activity := range activities {
		if err := logger.LogActivity(activity); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}

	fmt.Println("Activities logged successfully")
}