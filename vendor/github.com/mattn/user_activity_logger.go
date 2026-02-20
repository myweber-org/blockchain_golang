package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ActivityType string

const (
	Login      ActivityType = "LOGIN"
	Logout     ActivityType = "LOGOUT"
	ViewPage   ActivityType = "VIEW_PAGE"
	UpdateData ActivityType = "UPDATE_DATA"
)

type UserActivity struct {
	UserID    string       `json:"user_id"`
	Action    ActivityType `json:"action"`
	Timestamp time.Time    `json:"timestamp"`
	Details   string       `json:"details,omitempty"`
}

func LogActivity(userID string, action ActivityType, details string) UserActivity {
	activity := UserActivity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now().UTC(),
		Details:   details,
	}
	return activity
}

func SaveActivityToFile(activity UserActivity, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(activity)
}

func main() {
	activities := []UserActivity{
		LogActivity("user123", Login, "Successful authentication"),
		LogActivity("user123", ViewPage, "/dashboard"),
		LogActivity("user456", UpdateData, "Updated profile information"),
		LogActivity("user123", Logout, "Session terminated"),
	}

	for _, activity := range activities {
		err := SaveActivityToFile(activity, "user_activities.json")
		if err != nil {
			fmt.Printf("Failed to save activity: %v\n", err)
		} else {
			fmt.Printf("Logged: %s - %s\n", activity.UserID, activity.Action)
		}
	}
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
		"Method: %s | Path: %s | Duration: %v | RemoteAddr: %s",
		r.Method,
		r.URL.Path,
		duration,
		r.RemoteAddr,
	)
}