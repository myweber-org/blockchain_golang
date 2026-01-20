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

func NewActivityLog(userID, action, details string) *ActivityLog {
    return &ActivityLog{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
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
    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    logger := log.New(file, "", 0)

    LogActivity(logger, "user123", "LOGIN", "User logged in from web browser")
    LogActivity(logger, "user456", "UPLOAD", "File uploaded: report.pdf")
    LogActivity(logger, "user123", "LOGOUT", "Session terminated")

    fmt.Println("Activity logging completed. Check activity.log file.")
}