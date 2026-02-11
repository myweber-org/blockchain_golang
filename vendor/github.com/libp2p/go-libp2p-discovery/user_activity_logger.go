package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details"`
}

func NewActivityLog(userID, action, details string) *ActivityLog {
    return &ActivityLog{
        Timestamp: time.Now().UTC(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }
}

func (log *ActivityLog) SaveToFile(filename string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(log)
}

func main() {
    log := NewActivityLog("user123", "login", "Successful authentication via OAuth2")
    
    err := log.SaveToFile("activity_logs.json")
    if err != nil {
        fmt.Printf("Failed to save log: %v\n", err)
        return
    }
    
    fmt.Printf("Activity logged: %s performed %s at %s\n", 
        log.UserID, log.Action, log.Timestamp.Format(time.RFC3339))
}