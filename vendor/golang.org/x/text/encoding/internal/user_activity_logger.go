package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	userHits    map[string][]time.Time
	windowSize  time.Duration
	maxRequests int
}

func NewActivityLogger(window time.Duration, limit int) *ActivityLogger {
	return &ActivityLogger{
		userHits:    make(map[string][]time.Time),
		windowSize:  window,
		maxRequests: limit,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIP := r.RemoteAddr
		currentTime := time.Now()

		al.mu.Lock()
		defer al.mu.Unlock()

		al.cleanupOldHits(userIP, currentTime)

		if len(al.userHits[userIP]) >= al.maxRequests {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for IP: %s", userIP)
			return
		}

		al.userHits[userIP] = append(al.userHits[userIP], currentTime)
		log.Printf("Activity logged - IP: %s, Path: %s, Time: %v", userIP, r.URL.Path, currentTime)

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) cleanupOldHits(userIP string, currentTime time.Time) {
	hits := al.userHits[userIP]
	validHits := []time.Time{}

	for _, hitTime := range hits {
		if currentTime.Sub(hitTime) <= al.windowSize {
			validHits = append(validHits, hitTime)
		}
	}

	al.userHits[userIP] = validHits
}

func (al *ActivityLogger) GetUserActivity(userIP string) []time.Time {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return al.userHits[userIP]
}