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
	log.Printf("%s %s %d %s", r.Method, r.URL.Path, recorder.statusCode, duration)
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
	"context"
	"net/http"
	"sync"
	"time"
)

type ActivityEvent struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
	UserAgent string
}

type ActivityLogger struct {
	mu            sync.RWMutex
	events        []ActivityEvent
	rateLimiter   map[string]time.Time
	flushInterval time.Duration
	maxEvents     int
}

func NewActivityLogger(flushInterval time.Duration, maxEvents int) *ActivityLogger {
	return &ActivityLogger{
		events:        make([]ActivityEvent, 0, maxEvents),
		rateLimiter:   make(map[string]time.Time),
		flushInterval: flushInterval,
		maxEvents:     maxEvents,
	}
}

func (al *ActivityLogger) LogActivity(ctx context.Context, userID, action, ip, agent string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	al.mu.Lock()
	defer al.mu.Unlock()

	key := userID + ":" + action
	if lastTime, exists := al.rateLimiter[key]; exists {
		if time.Since(lastTime) < time.Minute {
			return nil
		}
	}

	event := ActivityEvent{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
		IPAddress: ip,
		UserAgent: agent,
	}

	al.events = append(al.events, event)
	al.rateLimiter[key] = event.Timestamp

	if len(al.events) >= al.maxEvents {
		go al.flushEvents()
	}

	return nil
}

func (al *ActivityLogger) flushEvents() {
	al.mu.Lock()
	if len(al.events) == 0 {
		al.mu.Unlock()
		return
	}

	events := make([]ActivityEvent, len(al.events))
	copy(events, al.events)
	al.events = al.events[:0]

	for k := range al.rateLimiter {
		delete(al.rateLimiter, k)
	}
	al.mu.Unlock()

	// In production, this would send to external service
	// For now just simulate processing
	time.Sleep(100 * time.Millisecond)
}

func (al *ActivityLogger) StartFlushWorker(ctx context.Context) {
	ticker := time.NewTicker(al.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			al.flushEvents()
			return
		case <-ticker.C:
			al.flushEvents()
		}
	}
}

func ActivityLoggingMiddleware(al *ActivityLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = "anonymous"
			}

			action := r.Method + " " + r.URL.Path
			ip := r.RemoteAddr
			agent := r.UserAgent()

			// Log asynchronously to not block request
			go al.LogActivity(ctx, userID, action, ip, agent)

			next.ServeHTTP(w, r)
		})
	}
}