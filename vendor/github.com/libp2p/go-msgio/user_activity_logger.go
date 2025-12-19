package main

import (
    "encoding/json"
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

func logActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now().UTC(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(logEntry); err != nil {
        return err
    }

    return nil
}

func main() {
    err := logActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100")
    if err != nil {
        log.Fatal("Failed to log activity:", err)
    }
    log.Println("Activity logged successfully")
}