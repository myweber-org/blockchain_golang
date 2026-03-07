
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
	startTime := time.Now()
	al.handler.ServeHTTP(w, r)
	duration := time.Since(startTime)

	log.Printf("Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}
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
}
package middleware

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	rateLimiter map[string][]time.Time
	window      time.Duration
	maxRequests int
}

func NewActivityLogger(window time.Duration, maxRequests int) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string][]time.Time),
		window:      window,
		maxRequests: maxRequests,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path
		method := r.Method

		if !al.allowRequest(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			select {
			case <-ctx.Done():
				log.Printf("Logging timeout for %s %s", method, path)
			default:
				log.Printf("Activity: %s %s from %s (UA: %s) took %v",
					method, path, clientIP, userAgent, duration)
			}
		}()
	})
}

func (al *ActivityLogger) allowRequest(clientIP string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	if requests, exists := al.rateLimiter[clientIP]; exists {
		var validRequests []time.Time
		for _, t := range requests {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}
		al.rateLimiter[clientIP] = validRequests

		if len(validRequests) >= al.maxRequests {
			return false
		}
	}

	al.rateLimiter[clientIP] = append(al.rateLimiter[clientIP], now)
	return true
}

func (al *ActivityLogger) Cleanup() {
	ticker := time.NewTicker(al.window * 2)
	defer ticker.Stop()

	for range ticker.C {
		al.mu.Lock()
		windowStart := time.Now().Add(-al.window)
		for ip, requests := range al.rateLimiter {
			var validRequests []time.Time
			for _, t := range requests {
				if t.After(windowStart) {
					validRequests = append(validRequests, t)
				}
			}
			if len(validRequests) == 0 {
				delete(al.rateLimiter, ip)
			} else {
				al.rateLimiter[ip] = validRequests
			}
		}
		al.mu.Unlock()
	}
}