package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	userLimits  map[string]time.Time
	rateLimit   time.Duration
}

func NewActivityLogger(limit time.Duration) *ActivityLogger {
	return &ActivityLogger{
		userLimits: make(map[string]time.Time),
		rateLimit:  limit,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		if !al.checkRateLimit(userID) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		al.logActivity(userID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) checkRateLimit(userID string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	lastTime, exists := al.userLimits[userID]
	now := time.Now()

	if exists && now.Sub(lastTime) < al.rateLimit {
		return false
	}

	al.userLimits[userID] = now
	return true
}

func (al *ActivityLogger) logActivity(userID, method, path string) {
	log.Printf("Activity: user=%s method=%s path=%s time=%s",
		userID, method, path, time.Now().Format(time.RFC3339))
}