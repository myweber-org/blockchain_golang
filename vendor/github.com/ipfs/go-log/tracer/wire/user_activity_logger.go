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

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(activity); err != nil {
        return err
    }

    fmt.Printf("Logged: %s performed %s at %s\n", userID, action, activity.Timestamp.Format(time.RFC3339))
    return nil
}

func main() {
    if err := logActivity("user123", "login", "Successful authentication"); err != nil {
        log.Fatal(err)
    }

    if err := logActivity("user456", "file_upload", "uploaded profile.jpg"); err != nil {
        log.Fatal(err)
    }

    if err := logActivity("user123", "logout", "Session terminated"); err != nil {
        log.Fatal(err)
    }
}