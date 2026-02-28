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

	log.Printf(
		"Method: %s | Path: %s | Duration: %v | Timestamp: %s",
		r.Method,
		r.URL.Path,
		duration,
		time.Now().Format(time.RFC3339),
	)
}package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	userLimits  map[string]*rateLimiter
	maxRequests int
	window      time.Duration
}

type rateLimiter struct {
	count    int
	lastSeen time.Time
}

func NewActivityLogger(maxRequests int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		userLimits:  make(map[string]*rateLimiter),
		maxRequests: maxRequests,
		window:      window,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		if !al.allowRequest(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("User %s accessed %s %s - Duration: %v",
			userID, r.Method, r.URL.Path, duration)
	})
}

func (al *ActivityLogger) allowRequest(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	limiter, exists := al.userLimits[userID]

	if !exists {
		al.userLimits[userID] = &rateLimiter{
			count:    1,
			lastSeen: now,
		}
		return true
	}

	if now.Sub(limiter.lastSeen) > al.window {
		limiter.count = 1
		limiter.lastSeen = now
		return true
	}

	if limiter.count >= al.maxRequests {
		return false
	}

	limiter.count++
	limiter.lastSeen = now
	return true
}

func (al *ActivityLogger) CleanupInactiveUsers(maxInactive time.Duration) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		al.mu.Lock()
		now := time.Now()
		for userID, limiter := range al.userLimits {
			if now.Sub(limiter.lastSeen) > maxInactive {
				delete(al.userLimits, userID)
			}
		}
		al.mu.Unlock()
	}
}