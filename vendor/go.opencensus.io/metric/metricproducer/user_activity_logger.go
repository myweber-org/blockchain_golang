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
	
	al.handler.ServeHTTP(w, r)
	
	duration := time.Since(start)
	
	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
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
	
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s",
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
    "log"
    "net/http"
    "time"
)

type ActivityLogger struct {
    handler http.Handler
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    al.handler.ServeHTTP(w, r)
    duration := time.Since(start)

    log.Printf("[%s] %s %s %s %v",
        time.Now().Format(time.RFC3339),
        r.Method,
        r.URL.Path,
        r.RemoteAddr,
        duration)
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
    return &ActivityLogger{handler: handler}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "ok"}`))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", apiHandler)

    wrappedMux := NewActivityLogger(mux)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}package main

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
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(activity); err != nil {
		return fmt.Errorf("failed to encode activity: %w", err)
	}

	return nil
}

func main() {
	if err := logActivity("user123", "login", "successful authentication"); err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	if err := logActivity("user456", "file_upload", "uploaded profile.jpg"); err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	fmt.Println("Activity logging completed")
}