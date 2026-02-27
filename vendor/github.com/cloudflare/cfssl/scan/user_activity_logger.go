package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		userAgent := r.UserAgent()
		clientIP := r.RemoteAddr

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		al.Logger.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | IP: %s | Agent: %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			clientIP,
			userAgent,
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityLog struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) ActivityLog {
    logEntry := ActivityLog{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    logData, err := json.MarshalIndent(logEntry, "", "  ")
    if err != nil {
        fmt.Printf("Error marshaling log: %v\n", err)
        return logEntry
    }

    fmt.Fprintf(os.Stdout, "%s\n", logData)
    return logEntry
}

func main() {
    logActivity("user_123", "login", "User logged in from web browser")
    logActivity("user_456", "purchase", "Purchased item: premium_subscription")
    logActivity("user_123", "logout", "Session terminated after 2 hours")
}