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
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		logEntry := ActivityLog{
			Timestamp: start.Format(time.RFC3339),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IP:        r.RemoteAddr,
		}
		
		logData, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("Failed to marshal log entry: %v", err)
		} else {
			log.Printf("Activity: %s", string(logData))
		}
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		log.Printf("Request completed in %v", duration)
	})
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request processed successfully")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	
	wrappedMux := loggingMiddleware(mux)
	
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", wrappedMux)
	if err != nil {
		log.Fatal(err)
	}
}