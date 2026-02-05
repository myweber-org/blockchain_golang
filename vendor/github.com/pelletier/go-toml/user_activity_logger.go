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
		now := time.Now()

		al.mu.Lock()
		defer al.mu.Unlock()

		hits := al.userHits[userIP]
		validHits := []time.Time{}
		for _, hit := range hits {
			if now.Sub(hit) <= al.windowSize {
				validHits = append(validHits, hit)
			}
		}

		if len(validHits) >= al.maxRequests {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for IP: %s", userIP)
			return
		}

		validHits = append(validHits, now)
		al.userHits[userIP] = validHits

		log.Printf("Activity from %s: %s %s", userIP, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(al.windowSize * 2)
	go func() {
		for range ticker.C {
			al.mu.Lock()
			now := time.Now()
			for ip, hits := range al.userHits {
				validHits := []time.Time{}
				for _, hit := range hits {
					if now.Sub(hit) <= al.windowSize {
						validHits = append(validHits, hit)
					}
				}
				if len(validHits) == 0 {
					delete(al.userHits, ip)
				} else {
					al.userHits[ip] = validHits
				}
			}
			al.mu.Unlock()
		}
	}()
}