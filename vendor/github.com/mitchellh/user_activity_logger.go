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
	
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
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
    "log"
    "net/http"
    "time"
)

type Activity struct {
    SessionID  string    `json:"session_id"`
    UserID     string    `json:"user_id"`
    Action     string    `json:"action"`
    Timestamp  time.Time `json:"timestamp"`
    UserAgent  string    `json:"user_agent"`
    IPAddress  string    `json:"ip_address"`
}

var activityLog []Activity

func logActivity(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var activity Activity
    if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    activity.Timestamp = time.Now()
    activity.IPAddress = r.RemoteAddr
    activity.UserAgent = r.UserAgent()

    activityLog = append(activityLog, activity)

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Activity logged successfully")
}

func getActivities(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(activityLog)
}

func main() {
    http.HandleFunc("/log", logActivity)
    http.HandleFunc("/activities", getActivities)

    log.Println("Server starting on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelWarn
	LevelError
)

type ActivityLogger struct {
	level    LogLevel
	format   string
	output   *log.Logger
	next     http.Handler
}

func NewActivityLogger(level LogLevel, format string, output *log.Logger, next http.Handler) *ActivityLogger {
	return &ActivityLogger{
		level:  level,
		format: format,
		output: output,
		next:   next,
	}
}

func (l *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	l.next.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	if l.shouldLog(recorder.statusCode) {
		l.logActivity(r, recorder.statusCode, duration)
	}
}

func (l *ActivityLogger) shouldLog(statusCode int) bool {
	if l.level == LevelError && statusCode < 400 {
		return false
	}
	if l.level == LevelWarn && statusCode < 300 {
		return false
	}
	return true
}

func (l *ActivityLogger) logActivity(r *http.Request, statusCode int, duration time.Duration) {
	switch l.format {
	case "json":
		l.output.Printf(`{"timestamp":"%s","method":"%s","path":"%s","status":%d,"duration_ms":%d}`,
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			statusCode,
			duration.Milliseconds())
	default:
		l.output.Printf("[%s] %s %s %d %v",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			statusCode,
			duration)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}