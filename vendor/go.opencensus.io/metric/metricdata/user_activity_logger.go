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
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | User-Agent: %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.UserAgent(),
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

type ActivityLog struct {
	UserID    string
	IPAddress string
	Method    string
	Path      string
	Timestamp time.Time
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			IPAddress: r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: start,
		}

		log.Printf("Activity: %s %s by %s from %s", 
			activity.Method, 
			activity.Path, 
			activity.UserID, 
			activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "sync"
    "time"
)

type ActivityEvent struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    SessionID string    `json:"session_id"`
    Metadata  string    `json:"metadata,omitempty"`
}

type ActivityLogger struct {
    mu       sync.Mutex
    file     *os.File
    encoder  *json.Encoder
    session  string
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    return &ActivityLogger{
        file:    file,
        encoder: json.NewEncoder(file),
        session: generateSessionID(),
    }, nil
}

func generateSessionID() string {
    return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

func (l *ActivityLogger) LogActivity(userID, action, metadata string) error {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    event := ActivityEvent{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        SessionID: l.session,
        Metadata:  metadata,
    }
    
    return l.encoder.Encode(event)
}

func (l *ActivityLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.file.Close()
}

func main() {
    logger, err := NewActivityLogger("activity.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    err = logger.LogActivity("user123", "login", "browser:chrome")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }
    
    err = logger.LogActivity("user123", "view_page", "page:/dashboard")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }
    
    err = logger.LogActivity("user123", "logout", "")
    if err != nil {
        log.Println("Failed to log activity:", err)
    }
    
    fmt.Println("Activity logging completed")
}