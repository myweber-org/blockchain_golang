package main

import (
    "fmt"
    "net/http"
    "time"
)

type Metrics struct {
    RequestCount    int
    TotalLatency   time.Duration
    ErrorCount     int
}

var metrics = &Metrics{}

func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
        
        next(recorder, r)
        
        duration := time.Since(start)
        metrics.RequestCount++
        metrics.TotalLatency += duration
        
        if recorder.statusCode >= 400 {
            metrics.ErrorCount++
        }
    }
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
    r.statusCode = code
    r.ResponseWriter.WriteHeader(code)
}

func handler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(10 * time.Millisecond)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Request processed")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "Internal server error")
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
    avgLatency := time.Duration(0)
    if metrics.RequestCount > 0 {
        avgLatency = metrics.TotalLatency / time.Duration(metrics.RequestCount)
    }
    
    errorRate := 0.0
    if metrics.RequestCount > 0 {
        errorRate = float64(metrics.ErrorCount) / float64(metrics.RequestCount) * 100
    }
    
    fmt.Fprintf(w, "Requests: %d\n", metrics.RequestCount)
    fmt.Fprintf(w, "Average Latency: %v\n", avgLatency)
    fmt.Fprintf(w, "Error Rate: %.2f%%\n", errorRate)
}

func main() {
    http.HandleFunc("/", metricsMiddleware(handler))
    http.HandleFunc("/error", metricsMiddleware(errorHandler))
    http.HandleFunc("/metrics", metricsHandler)
    
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}