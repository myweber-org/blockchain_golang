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
		userAgent := r.UserAgent()
		clientIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		al.Logger.Printf(
			"User activity: IP=%s Method=%s Path=%s Agent=%s Duration=%v",
			clientIP,
			method,
			path,
			userAgent,
			duration,
		)
	})
}