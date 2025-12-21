package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type Activity struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
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

func (l *ActivityLogger) LogActivity(userID, action, details string) error {
    activity := Activity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

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

    activities := []struct {
        userID  string
        action  string
        details string
    }{
        {"user_123", "login", "Successful authentication"},
        {"user_456", "purchase", "Order ID: ORD-78910"},
        {"user_123", "logout", "Session duration: 45m"},
        {"user_789", "view_page", "/products/electronics"},
    }

    for _, act := range activities {
        if err := logger.LogActivity(act.userID, act.action, act.details); err != nil {
            log.Printf("Failed to log activity: %v", err)
        }
    }

    fmt.Println("Activity logging completed")
}