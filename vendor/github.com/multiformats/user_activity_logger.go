package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityEvent struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
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

func (l *ActivityLogger) LogActivity(userID, eventType, details string) error {
    event := ActivityEvent{
        UserID:    userID,
        EventType: eventType,
        Timestamp: time.Now(),
        Details:   details,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    data = append(data, '\n')
    _, err = l.logFile.Write(data)
    return err
}

func (l *ActivityLogger) Close() error {
    return l.logFile.Close()
}

func main() {
    logger, err := NewActivityLogger("activity.log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    err = logger.LogActivity("user123", "login", "User logged in from web browser")
    if err != nil {
        log.Fatal(err)
    }

    err = logger.LogActivity("user123", "search", "Searched for 'golang tutorials'")
    if err != nil {
        log.Fatal(err)
    }

    err = logger.LogActivity("user456", "logout", "Session expired")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Activity logged successfully")
}package middleware

import (
    "log"
    "net/http"
    "time"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func ActivityLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(rw, r)
        
        duration := time.Since(start)
        log.Printf("[%s] %s %s %d %v",
            time.Now().Format("2006-01-02 15:04:05"),
            r.Method,
            r.URL.Path,
            rw.statusCode,
            duration,
        )
    })
}package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    IPAddress string    `json:"ip_address"`
}

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

func (rl *RateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-rl.window)

    timestamps := rl.requests[ip]
    validRequests := []time.Time{}

    for _, ts := range timestamps {
        if ts.After(windowStart) {
            validRequests = append(validRequests, ts)
        }
    }

    if len(validRequests) >= rl.limit {
        return false
    }

    validRequests = append(validRequests, now)
    rl.requests[ip] = validRequests

    return true
}

func loggingMiddleware(next http.Handler) http.Handler {
    limiter := NewRateLimiter(100, time.Minute)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr

        if !limiter.Allow(ip) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }

        logEntry := ActivityLog{
            Timestamp: time.Now(),
            UserID:    userID,
            Action:    r.Method + " " + r.URL.Path,
            IPAddress: ip,
        }

        logData, err := json.Marshal(logEntry)
        if err != nil {
            log.Printf("Failed to marshal log entry: %v", err)
        } else {
            log.Printf("Activity: %s", string(logData))
        }

        next.ServeHTTP(w, r)
    })
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "status":  "success",
        "message": "Request processed",
        "time":    time.Now().Format(time.RFC3339),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api", apiHandler)

    wrappedMux := loggingMiddleware(mux)

    server := &http.Server{
        Addr:         ":8080",
        Handler:      wrappedMux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    log.Println("Server starting on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}