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
}package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityType string

const (
	Login    ActivityType = "LOGIN"
	Logout   ActivityType = "LOGOUT"
	Purchase ActivityType = "PURCHASE"
	View     ActivityType = "VIEW"
)

type UserActivity struct {
	UserID    string
	Action    ActivityType
	Timestamp time.Time
	Details   string
}

var activityLog []UserActivity

func logActivity(userID string, action ActivityType, details string) {
	activity := UserActivity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
		Details:   details,
	}
	activityLog = append(activityLog, activity)
	fmt.Printf("Logged: %s - %s at %s\n", userID, action, activity.Timestamp.Format(time.RFC3339))
}

func saveLogToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, activity := range activityLog {
		line := fmt.Sprintf("%s,%s,%s,%s\n",
			activity.UserID,
			activity.Action,
			activity.Timestamp.Format(time.RFC3339),
			activity.Details)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	logActivity("user123", Login, "Successful login from IP 192.168.1.100")
	logActivity("user123", View, "Viewed product page: product_456")
	logActivity("user456", Purchase, "Purchased item: premium_subscription")
	logActivity("user123", Logout, "Session terminated")

	err := saveLogToFile("user_activities.csv")
	if err != nil {
		log.Fatal("Failed to save log:", err)
	}
	fmt.Println("Activity log saved to user_activities.csv")
}