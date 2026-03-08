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
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	logger     *log.Logger
	limiter    *rate.Limiter
	ignoredPaths map[string]bool
}

func NewActivityLogger(logger *log.Logger, rps float64, burst int) *ActivityLogger {
	return &ActivityLogger{
		logger:     logger,
		limiter:    rate.NewLimiter(rate.Limit(rps), burst),
		ignoredPaths: make(map[string]bool),
	}
}

func (al *ActivityLogger) IgnorePath(path string) {
	al.ignoredPaths[path] = true
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if al.ignoredPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		if !al.limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		start := time.Now()
		userAgent := r.Header.Get("User-Agent")
		clientIP := r.RemoteAddr

		al.logger.Printf("Request started: %s %s from %s (UA: %s)", 
			r.Method, r.URL.Path, clientIP, userAgent)

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(ctx))

		duration := time.Since(start)
		al.logger.Printf("Request completed: %s %s - Status: %d, Duration: %v", 
			r.Method, r.URL.Path, rw.statusCode, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}