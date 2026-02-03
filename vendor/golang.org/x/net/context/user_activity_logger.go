package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter *RateLimiter
}

type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	requests := rl.requests[ip]

	var validRequests []time.Time
	for _, t := range requests {
		if now.Sub(t) <= rl.window {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= rl.limit {
		return false
	}

	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests
	return true
}

func NewActivityLogger() *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: NewRateLimiter(100, time.Minute),
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path
		method := r.Method

		if !al.rateLimiter.Allow(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		defer func() {
			duration := time.Since(start)
			status := recorder.statusCode

			log.Printf("IP: %s | Method: %s | Path: %s | Status: %d | Duration: %v | Agent: %s",
				ip, method, path, status, duration, userAgent)

			if status >= 500 {
				log.Printf("ERROR: Server error occurred for IP %s on path %s", ip, path)
			}
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
	userAgent := r.UserAgent()
	clientIP := r.RemoteAddr
	requestPath := r.URL.Path
	requestMethod := r.Method

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("Activity: %s %s | IP: %s | Agent: %s | Duration: %v",
		requestMethod,
		requestPath,
		clientIP,
		userAgent,
		duration,
	)
}