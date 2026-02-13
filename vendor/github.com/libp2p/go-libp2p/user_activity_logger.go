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
	recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	log.Printf("[%s] %s %s %d %v",
		r.RemoteAddr,
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
	)
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
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	activities  map[string][]time.Time
	rateLimit   int
	window      time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		activities: make(map[string][]time.Time),
		rateLimit:  limit,
		window:     window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = r.RemoteAddr
		}

		if !al.checkRateLimit(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("User %s accessed %s at %s (took %v)", 
			userID, r.URL.Path, start.Format(time.RFC3339), duration)
	})
}

func (al *ActivityLogger) checkRateLimit(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	// Clean old entries
	validTimes := []time.Time{}
	for _, t := range al.activities[userID] {
		if t.After(windowStart) {
			validTimes = append(validTimes, t)
		}
	}

	if len(validTimes) >= al.rateLimit {
		return false
	}

	validTimes = append(validTimes, now)
	al.activities[userID] = validTimes
	return true
}

func (al *ActivityLogger) GetUserActivity(userID string) []time.Time {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return append([]time.Time{}, al.activities[userID]...)
}