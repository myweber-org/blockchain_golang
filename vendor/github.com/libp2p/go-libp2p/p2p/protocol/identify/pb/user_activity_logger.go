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
	ipAddress := r.RemoteAddr

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s from IP %s accessed %s %s - Duration: %v", userID, ipAddress, r.Method, r.URL.Path, duration)
}

func extractUserID(r *http.Request) string {
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return userID
	}
	return "anonymous"
}package middleware

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
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityType string

const (
    Login    ActivityType = "LOGIN"
    Logout   ActivityType = "LOGOUT"
    Purchase ActivityType = "PURCHASE"
    View     ActivityType = "VIEW"
)

type UserActivity struct {
    UserID    string       `json:"user_id"`
    Timestamp time.Time    `json:"timestamp"`
    Action    ActivityType `json:"action"`
    Details   string       `json:"details,omitempty"`
}

type ActivityLogger struct {
    logFile *os.File
    encoder *json.Encoder
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    return &ActivityLogger{
        logFile: file,
        encoder: json.NewEncoder(file),
    }, nil
}

func (l *ActivityLogger) LogActivity(userID string, action ActivityType, details string) error {
    activity := UserActivity{
        UserID:    userID,
        Timestamp: time.Now().UTC(),
        Action:    action,
        Details:   details,
    }
    return l.encoder.Encode(activity)
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
        action  ActivityType
        details string
    }{
        {"user123", Login, "Successful authentication"},
        {"user456", View, "Product page: laptop-x100"},
        {"user123", Purchase, "Order ID: ORD-78901"},
        {"user789", Logout, "Session ended"},
    }

    for _, act := range activities {
        if err := logger.LogActivity(act.userID, act.action, act.details); err != nil {
            fmt.Printf("Failed to log activity: %v\n", err)
        }
    }

    fmt.Println("Activity logging completed")
}