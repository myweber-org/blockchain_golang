package middleware

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
			"[%s] %s %s %d %v",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	})
}
package middleware

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
	Status    int
}

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(startTime)

	activity := ActivityLog{
		Timestamp: startTime,
		Method:    r.Method,
		Path:      r.URL.Path,
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
		Duration:  duration,
		Status:    recorder.statusCode,
	}

	al.logActivity(activity)
}

func (al *ActivityLogger) logActivity(activity ActivityLog) {
	log.Printf("ACTIVITY: %s | %s %s | %s | %s | %d | %v",
		activity.Timestamp.Format(time.RFC3339),
		activity.Method,
		activity.Path,
		activity.IPAddress,
		activity.UserAgent,
		activity.Status,
		activity.Duration,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}