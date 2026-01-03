package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Activity struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Endpoint  string    `json:"endpoint"`
	IPAddress string    `json:"ip_address"`
}

type ActivityLogger struct {
	mu       sync.RWMutex
	activities []Activity
	rateLimit map[string]time.Time
}

func NewActivityLogger() *ActivityLogger {
	return &ActivityLogger{
		activities: make([]Activity, 0),
		rateLimit:  make(map[string]time.Time),
	}
}

func (al *ActivityLogger) LogActivity(userID, action, endpoint, ip string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	key := ip + ":" + endpoint
	if lastTime, exists := al.rateLimit[key]; exists {
		if time.Since(lastTime) < time.Second {
			return false
		}
	}

	activity := Activity{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Endpoint:  endpoint,
		IPAddress: ip,
	}

	al.activities = append(al.activities, activity)
	al.rateLimit[key] = activity.Timestamp
	return true
}

func (al *ActivityLogger) GetActivities() []Activity {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return al.activities
}

func loggingMiddleware(al *ActivityLogger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		action := "ACCESS"
		if r.Method != "GET" {
			action = "MODIFY"
		}

		logged := al.LogActivity(userID, action, r.URL.Path, ip)
		if !logged {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func activityHandler(al *ActivityLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activities := al.GetActivities()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(activities)
	}
}

func main() {
	logger := NewActivityLogger()

	http.HandleFunc("/api/data", loggingMiddleware(logger, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Data endpoint response"))
	}))

	http.HandleFunc("/admin/activities", activityHandler(logger))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}