package middleware

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
		
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(lrw, r)
		
		duration := time.Since(start)
		
		activity := ActivityLog{
			Timestamp:  time.Now(),
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
}