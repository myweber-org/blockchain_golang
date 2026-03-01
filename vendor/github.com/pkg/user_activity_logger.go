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
		
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		
		al.Logger.Printf(
			"Method: %s | Path: %s | Status: %d | Duration: %v | User-Agent: %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
			r.UserAgent(),
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
	
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
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
		"Activity: %s %s | Duration: %v | RemoteAddr: %s | UserAgent: %s",
		r.Method,
		r.URL.Path,
		duration,
		r.RemoteAddr,
		r.UserAgent(),
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
	writer := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	
	al.handler.ServeHTTP(writer, r)
	
	duration := time.Since(start)
	log.Printf("%s %s %d %v", r.Method, r.URL.Path, writer.statusCode, duration)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	redisClient *redis.Client
	limiter     *rate.Limiter
	serviceName string
}

type ActivityEvent struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Timestamp time.Time `json:"timestamp"`
	Metadata  string    `json:"metadata,omitempty"`
}

func NewActivityLogger(redisAddr string, serviceName string, rps int) (*ActivityLogger, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &ActivityLogger{
		redisClient: rdb,
		limiter:     rate.NewLimiter(rate.Limit(rps), rps),
		serviceName: serviceName,
	}, nil
}

func (al *ActivityLogger) LogActivity(ctx context.Context, userID, action, resource, metadata string) error {
	if !al.limiter.Allow() {
		return fmt.Errorf("rate limit exceeded for activity logging")
	}

	event := ActivityEvent{
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Timestamp: time.Now().UTC(),
		Metadata:  metadata,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal activity event: %w", err)
	}

	key := fmt.Sprintf("activity:%s:%s:%d", al.serviceName, userID, time.Now().UnixNano())
	expiration := 24 * time.Hour

	if err := al.redisClient.Set(ctx, key, eventJSON, expiration).Err(); err != nil {
		return fmt.Errorf("failed to store activity in Redis: %w", err)
	}

	streamKey := fmt.Sprintf("activity_stream:%s", al.serviceName)
	streamValues := map[string]interface{}{
		"user_id":   userID,
		"action":    action,
		"resource":  resource,
		"metadata":  metadata,
		"timestamp": event.Timestamp.Format(time.RFC3339),
	}

	if err := al.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: streamValues,
	}).Err(); err != nil {
		return fmt.Errorf("failed to add to activity stream: %w", err)
	}

	return nil
}

func (al *ActivityLogger) GetUserActivities(ctx context.Context, userID string, limit int64) ([]ActivityEvent, error) {
	pattern := fmt.Sprintf("activity:%s:%s:*", al.serviceName, userID)
	keys, err := al.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get activity keys: %w", err)
	}

	if len(keys) == 0 {
		return []ActivityEvent{}, nil
	}

	if limit > 0 && int64(len(keys)) > limit {
		keys = keys[:limit]
	}

	var activities []ActivityEvent
	for _, key := range keys {
		val, err := al.redisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var event ActivityEvent
		if err := json.Unmarshal([]byte(val), &event); err != nil {
			continue
		}
		activities = append(activities, event)
	}

	return activities, nil
}

func (al *ActivityLogger) Close() error {
	return al.redisClient.Close()
}