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
}package middleware

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
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.RemoteAddr,
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

	log.Printf(
		"User activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}