package middleware

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	limiter *rate.Limiter
	logger  *log.Logger
}

func NewActivityLogger(rps int, logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{
		limiter: rate.NewLimiter(rate.Limit(rps), rps*2),
		logger:  logger,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !al.limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		defer func() {
			duration := time.Since(start)
			al.logger.Printf(
				"method=%s path=%s status=%d duration=%s remote=%s",
				r.Method,
				r.URL.Path,
				recorder.statusCode,
				duration,
				r.RemoteAddr,
			)
		}()

		next.ServeHTTP(recorder, r)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}