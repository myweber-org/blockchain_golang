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

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		al.Logger.Printf(
			"IP: %s | Method: %s | Path: %s | Status: %d | Duration: %v | Agent: %s",
			clientIP,
			method,
			path,
			recorder.statusCode,
			duration,
			userAgent,
		)
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
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"context"
	"net/http"
	"time"
)

type ActivityLogger struct {
	store      ActivityStore
	rateLimiter RateLimiter
}

type ActivityStore interface {
	LogActivity(ctx context.Context, userID string, action string, metadata map[string]interface{}) error
}

type RateLimiter interface {
	Allow(userID string) bool
}

func NewActivityLogger(store ActivityStore, limiter RateLimiter) *ActivityLogger {
	return &ActivityLogger{
		store:      store,
		rateLimiter: limiter,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		userID := extractUserID(r)
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !al.rateLimiter.Allow(userID) {
			next.ServeHTTP(w, r)
			return
		}

		action := r.Method + " " + r.URL.Path
		metadata := map[string]interface{}{
			"user_agent": r.UserAgent(),
			"ip_address": r.RemoteAddr,
			"timestamp":  time.Now().UTC(),
		}

		go func() {
			if err := al.store.LogActivity(ctx, userID, action, metadata); err != nil {
				logError(ctx, err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return parseToken(auth)
	}
	return ""
}

func parseToken(token string) string {
	return token
}

func logError(ctx context.Context, err error) {
}