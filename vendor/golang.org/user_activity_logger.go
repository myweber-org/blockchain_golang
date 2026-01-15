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
			"%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu      sync.RWMutex
	entries map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		entries: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (al *ActivityLogger) LogActivity(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	activities := al.entries[userID]
	var validActivities []time.Time
	for _, t := range activities {
		if t.After(windowStart) {
			validActivities = append(validActivities, t)
		}
	}

	if len(validActivities) >= al.limit {
		return false
	}

	validActivities = append(validActivities, now)
	al.entries[userID] = validActivities
	return true
}

func (al *ActivityLogger) Cleanup() {
	al.mu.Lock()
	defer al.mu.Unlock()

	windowStart := time.Now().Add(-al.window)
	for userID, activities := range al.entries {
		var validActivities []time.Time
		for _, t := range activities {
			if t.After(windowStart) {
				validActivities = append(validActivities, t)
			}
		}
		if len(validActivities) == 0 {
			delete(al.entries, userID)
		} else {
			al.entries[userID] = validActivities
		}
	}
}

func loggingMiddleware(al *ActivityLogger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !al.LogActivity(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		log.Printf("Activity logged for user %s: %s %s", userID, r.Method, r.URL.Path)
		next(w, r)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success","message":"Request processed"}`))
}

func main() {
	logger := NewActivityLogger(10, time.Minute*5)
	go func() {
		ticker := time.NewTicker(time.Hour)
		for range ticker.C {
			logger.Cleanup()
		}
	}()

	http.HandleFunc("/api", loggingMiddleware(logger, apiHandler))
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}