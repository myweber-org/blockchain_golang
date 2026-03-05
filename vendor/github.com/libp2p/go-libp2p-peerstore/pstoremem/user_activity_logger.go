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
	writer := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	al.handler.ServeHTTP(writer, r)

	duration := time.Since(start)
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		writer.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
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
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	requests := rl.requests[ip]
	validRequests := []time.Time{}
	for _, t := range requests {
		if t.After(windowStart) {
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

func ActivityLogger(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if !limiter.Allow(ip) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			log.Printf("IP: %s | Method: %s | Path: %s | Duration: %v", ip, r.Method, r.URL.Path, duration)
		})
	}
}

func main() {
	limiter := NewRateLimiter(10, time.Minute)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request logged"))
	})

	handler := ActivityLogger(limiter)(mux)
	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", handler)
}