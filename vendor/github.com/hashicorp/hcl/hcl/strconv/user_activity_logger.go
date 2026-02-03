package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type ActivityKey string

const UserActivityKey ActivityKey = "user_activity"

type UserActivity struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
	UserAgent string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := UserActivity{
			UserID:    extractUserID(r),
			Action:    r.Method + " " + r.URL.Path,
			Timestamp: time.Now().UTC(),
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		ctx := context.WithValue(r.Context(), UserActivityKey, activity)
		next.ServeHTTP(w, r.WithContext(ctx))

		go logActivity(activity)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return "authenticated_user"
	}
	return "anonymous"
}

func logActivity(activity UserActivity) {
	log.Printf("Activity: User=%s Action=%s IP=%s Time=%s",
		activity.UserID,
		activity.Action,
		activity.IPAddress,
		activity.Timestamp.Format(time.RFC3339))
}

func GetActivityFromContext(ctx context.Context) (UserActivity, bool) {
	activity, ok := ctx.Value(UserActivityKey).(UserActivity)
	return activity, ok
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp  time.Time
	Method     string
	Path       string
	RemoteAddr string
	UserAgent  string
	StatusCode int
	Duration   time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(lrw, r)
		
		duration := time.Since(start)
		activity := ActivityLog{
			Timestamp:  start,
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			StatusCode: lrw.statusCode,
			Duration:   duration,
		}
		
		log.Printf("ACTIVITY: %s %s %d %s %s",
			activity.Method,
			activity.Path,
			activity.StatusCode,
			activity.Duration,
			activity.RemoteAddr,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
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
    if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
        log.Fatal(err)
    }
}