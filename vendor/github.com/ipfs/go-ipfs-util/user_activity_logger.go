package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityEvent struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Details   string    `json:"details"`
}

type ActivityLogger struct {
	logFile *os.File
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{logFile: file}, nil
}

func (al *ActivityLogger) LogActivity(userID, eventType, details string) error {
	event := ActivityEvent{
		Timestamp: time.Now(),
		UserID:    userID,
		EventType: eventType,
		Details:   details,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	eventJSON = append(eventJSON, '\n')
	_, err = al.logFile.Write(eventJSON)
	return err
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("user_activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	err = logger.LogActivity("user123", "LOGIN", "User logged in from IP 192.168.1.100")
	if err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	err = logger.LogActivity("user123", "VIEW_PAGE", "Accessed dashboard page")
	if err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	fmt.Println("Activity logging completed")
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
		"[%s] %s %s %s %v",
		time.Now().Format("2006-01-02 15:04:05"),
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}