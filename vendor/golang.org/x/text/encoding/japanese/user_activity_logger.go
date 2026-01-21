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

func NewActivityLog(userID, action, resource string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}
}

func (al *ActivityLog) SaveToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(al)
	if err != nil {
		return err
	}

	_, err = file.Write(append(data, '\n'))
	return err
}

func main() {
	logger := NewActivityLog("user_12345", "CREATE", "/api/v1/documents")
	
	err := logger.SaveToFile("activity_logs.json")
	if err != nil {
		log.Fatal("Failed to save activity log:", err)
	}
	
	fmt.Printf("Activity logged: %s performed %s on %s at %s\n",
		logger.UserID,
		logger.Action,
		logger.Resource,
		logger.Timestamp.Format(time.RFC3339))
}