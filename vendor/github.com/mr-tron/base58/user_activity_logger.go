package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityEvent struct {
	Timestamp time.Time
	UserID    string
	EventType string
	Details   string
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

func (al *ActivityLogger) LogActivity(userID, eventType, details string) {
	event := ActivityEvent{
		Timestamp: time.Now(),
		UserID:    userID,
		EventType: eventType,
		Details:   details,
	}

	logEntry := fmt.Sprintf("%s | User: %s | Event: %s | Details: %s\n",
		event.Timestamp.Format("2006-01-02 15:04:05"),
		event.UserID,
		event.EventType,
		event.Details)

	if _, err := al.logFile.WriteString(logEntry); err != nil {
		log.Printf("Failed to write activity log: %v", err)
	}
}

func (al *ActivityLogger) Close() {
	if al.logFile != nil {
		al.logFile.Close()
	}
}

func main() {
	logger, err := NewActivityLogger("user_activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	logger.LogActivity("user123", "LOGIN", "Successful authentication")
	logger.LogActivity("user123", "VIEW_PAGE", "Accessed dashboard")
	logger.LogActivity("user456", "UPLOAD", "File: report.pdf")
	logger.LogActivity("user123", "LOGOUT", "Session ended")

	fmt.Println("Activity logging completed. Check user_activity.log for details.")
}