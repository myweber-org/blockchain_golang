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
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type Activity struct {
    SessionID  string    `json:"session_id"`
    UserID     string    `json:"user_id"`
    Action     string    `json:"action"`
    Timestamp  time.Time `json:"timestamp"`
    UserAgent  string    `json:"user_agent"`
    IPAddress  string    `json:"ip_address"`
}

var activityLog []Activity

func logActivity(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var activity Activity
    if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    activity.Timestamp = time.Now()
    activity.IPAddress = r.RemoteAddr
    activity.UserAgent = r.UserAgent()

    activityLog = append(activityLog, activity)

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Activity logged successfully")
}

func getActivities(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(activityLog)
}

func main() {
    http.HandleFunc("/log", logActivity)
    http.HandleFunc("/activities", getActivities)

    log.Println("Server starting on port 8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}