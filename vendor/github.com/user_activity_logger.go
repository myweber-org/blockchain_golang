package middleware

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

		log.Printf("Activity: %s %s by %s from %s", activity.Method, activity.Path, activity.UserID, activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
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
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("API response"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)
    
    loggedMux := NewActivityLogger(mux)
    
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", loggedMux)
}