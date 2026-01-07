package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type UserActivity struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) error {
    activity := UserActivity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    file, err := os.OpenFile("user_activities.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
    if err := logActivity("user123", "login", "Successful authentication"); err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    if err := logActivity("user456", "file_upload", "Uploaded profile picture"); err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    fmt.Println("Activity logging completed")
}