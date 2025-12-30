
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

func NewActivityLog(userID, action, details string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
}

func (al *ActivityLog) ToJSON() ([]byte, error) {
	return json.MarshalIndent(al, "", "  ")
}

func LogActivity(logger *log.Logger, userID, action, details string) {
	activity := NewActivityLog(userID, action, details)
	jsonData, err := activity.ToJSON()
	if err != nil {
		logger.Printf("Failed to marshal activity log: %v", err)
		return
	}
	logger.Println(string(jsonData))
}

func main() {
	logFile, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	activityLogger := log.New(logFile, "", 0)

	LogActivity(activityLogger, "user123", "LOGIN", "User logged in from web browser")
	LogActivity(activityLogger, "user456", "UPDATE_PROFILE", "Changed email address")
	LogActivity(activityLogger, "user789", "LOGOUT", "Session expired")

	fmt.Println("Activity logging completed. Check activity.log file.")
}