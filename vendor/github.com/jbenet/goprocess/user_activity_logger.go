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
	log.Printf("User %s from %s accessed %s %s - Duration: %v",
		userID, ipAddress, r.Method, r.URL.Path, duration)
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
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Endpoint  string
	Method    string
	Timestamp time.Time
	IPAddress string
}

var activityChannel = make(chan ActivityLog, 100)

func activityLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		activity := ActivityLog{
			UserID:    extractUserID(r),
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: r.RemoteAddr,
		}

		select {
		case activityChannel <- activity:
		default:
			log.Println("Activity channel full, dropping log entry")
		}
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return auth[:8]
	}
	return "anonymous"
}

func processActivityLogs() {
	for activity := range activityChannel {
		log.Printf("Activity: %s %s by %s from %s at %s",
			activity.Method,
			activity.Endpoint,
			activity.UserID,
			activity.IPAddress,
			activity.Timestamp.Format(time.RFC3339))
	}
}

func main() {
	go processActivityLogs()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Users endpoint"))
	})
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Orders endpoint"))
	})

	handler := activityLoggerMiddleware(mux)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}