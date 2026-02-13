package middleware

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	limiter *rate.Limiter
}

func NewActivityLogger(rps int) *ActivityLogger {
	return &ActivityLogger{
		limiter: rate.NewLimiter(rate.Limit(rps), rps*2),
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !al.limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			duration := time.Since(start)
			al.logActivity(r, recorder.statusCode, duration)
		}()

		next.ServeHTTP(recorder, r)
	})
}

func (al *ActivityLogger) logActivity(r *http.Request, status int, duration time.Duration) {
	log.Printf("ACTIVITY: %s %s %d %s %s",
		r.Method,
		r.URL.Path,
		status,
		duration.Round(time.Millisecond),
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
}