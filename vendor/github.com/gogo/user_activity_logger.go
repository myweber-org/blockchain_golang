package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityLog struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) {
    activity := ActivityLog{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
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
    logActivity("user123", "login", "User logged in from web browser")
    logActivity("user456", "purchase", "Purchased item ID: 789")
    logActivity("user123", "logout", "Session ended normally")

    fmt.Println("Activity logging completed. Check activity.log file.")
}