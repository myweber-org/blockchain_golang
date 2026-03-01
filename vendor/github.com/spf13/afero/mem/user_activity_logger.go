package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter map[string]time.Time
	window      time.Duration
}

func NewActivityLogger(window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string]time.Time),
		window:      window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		now := time.Now()

		if lastSeen, exists := al.rateLimiter[clientIP]; exists {
			if now.Sub(lastSeen) < al.window {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}

		al.rateLimiter[clientIP] = now

		log.Printf("Activity: %s %s from %s", r.Method, r.URL.Path, clientIP)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for ip, lastSeen := range al.rateLimiter {
			if now.Sub(lastSeen) > 24*time.Hour {
				delete(al.rateLimiter, ip)
			}
		}
	}
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
	logActivity("user123", "login", "User logged in from web browser")
	logActivity("user456", "purchase", "Purchased item ID: 789")
	logActivity("user123", "logout", "Session ended after 30 minutes")

	fmt.Println("Activity logging completed. Check activity.log for details.")
}package middleware

import (
    "log"
    "net/http"
    "time"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func ActivityLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(rw, r)

        duration := time.Since(start)
        log.Printf(
            "%s %s %d %s %s",
            r.Method,
            r.URL.Path,
            rw.statusCode,
            duration,
            r.RemoteAddr,
        )
    })
}