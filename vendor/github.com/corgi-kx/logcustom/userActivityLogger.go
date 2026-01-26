package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	UserID    string
	Method    string
	Path      string
	Status    int
	Duration  time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		userID := extractUserID(r)
		
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		
		activity := ActivityLog{
			Timestamp: time.Now(),
			UserID:    userID,
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    lrw.statusCode,
			Duration:  time.Since(start),
		}
		
		logActivity(activity)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return hashString(auth[:min(16, len(auth))])
	}
	return "anonymous"
}

func hashString(s string) string {
	hash := 0
	for _, ch := range s {
		hash = 31*hash + int(ch)
	}
	return string(rune(hash % 26 + 97))
}

func logActivity(activity ActivityLog) {
	log.Printf("[ACTIVITY] %s | User: %s | %s %s -> %d (%v)",
		activity.Timestamp.Format("2006-01-02 15:04:05"),
		activity.UserID,
		activity.Method,
		activity.Path,
		activity.Status,
		activity.Duration,
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}