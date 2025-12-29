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

func logActivity(userID, action, details string) {
    activity := ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    logEntry, err := json.Marshal(activity)
    if err != nil {
        log.Printf("Failed to marshal activity log: %v", err)
        return
    }

    logFile, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Failed to open log file: %v", err)
        return
    }
    defer logFile.Close()

    if _, err := logFile.Write(append(logEntry, '\n')); err != nil {
        log.Printf("Failed to write log entry: %v", err)
    }
}

func main() {
    logActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100")
    logActivity("user456", "PURCHASE", "Purchased item ID: ITM-789")
    logActivity("user123", "LOGOUT", "Session duration: 15m30s")

    fmt.Println("Activity logging completed. Check activity.log file.")
}