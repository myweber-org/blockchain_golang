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
	userAgent := r.UserAgent()
	clientIP := r.RemoteAddr
	requestPath := r.URL.Path

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("Activity: %s | %s | %s | %v", clientIP, userAgent, requestPath, duration)
}package main

import (
	"encoding/json"
	"fmt"
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

func logActivity(userID, action, details string) {
	activity := ActivityLog{
		Timestamp: time.Now(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(activity); err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
}

func main() {
	logActivity("user123", "LOGIN", "User logged in from web browser")
	logActivity("user456", "UPDATE_PROFILE", "Changed email address")
	logActivity("user123", "LOGOUT", "Session expired after 30 minutes")

	fmt.Println("Activity logging completed. Check activity.log for details.")
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
	log.Printf("%s %s %s %v", r.RemoteAddr, r.Method, r.URL.Path, duration)
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", apiHandler)
	wrappedMux := NewActivityLogger(mux)
	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", wrappedMux)
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IPAddress string
	Duration  time.Duration
}

type ActivityLogger struct {
	activities chan ActivityLog
}

func NewActivityLogger(bufferSize int) *ActivityLogger {
	al := &ActivityLogger{
		activities: make(chan ActivityLog, bufferSize),
	}
	go al.processLogs()
	return al
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next.ServeHTTP(w, r)
		
		activity := ActivityLog{
			Timestamp: time.Now(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
			Duration:  time.Since(start),
		}
		
		select {
		case al.activities <- activity:
		default:
			log.Println("Activity log buffer full, dropping entry")
		}
	})
}

func (al *ActivityLogger) processLogs() {
	for activity := range al.activities {
		log.Printf("Activity: %s %s from %s (UA: %s) took %v",
			activity.Method,
			activity.Path,
			activity.IPAddress,
			activity.UserAgent,
			activity.Duration,
		)
	}
}

func (al *ActivityLogger) Close() {
	close(al.activities)
}