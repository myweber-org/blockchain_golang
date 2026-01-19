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
}package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Activity struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
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

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
	activity := Activity{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	data, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = al.logFile.Write(data)
	return err
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("user_activities.json")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	err = logger.LogActivity("user123", "login", "Successful authentication")
	if err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	err = logger.LogActivity("user123", "file_upload", "Uploaded document.pdf")
	if err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	fmt.Println("Activities logged successfully")
}