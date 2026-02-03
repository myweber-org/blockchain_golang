package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userID := extractUserID(r)
	ip := r.RemoteAddr
	method := r.Method
	path := r.URL.Path

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s from %s %s %s completed in %v", userID, ip, method, path, duration)
}

func extractUserID(r *http.Request) string {
	if user := r.Header.Get("X-User-ID"); user != "" {
		return user
	}
	return "anonymous"
}package main

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
    Details   string    `json:"details"`
}

func logActivity(userID, action, details string) error {
    logEntry := ActivityLog{
        Timestamp: time.Now(),
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
    return encoder.Encode(logEntry)
}

func main() {
    err := logActivity("user123", "login", "Successful authentication")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = logActivity("user123", "upload", "File: report.pdf")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    log.Println("Activity logging completed")
}