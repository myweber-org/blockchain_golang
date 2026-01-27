package main

import (
    "fmt"
    "net/http"
    "time"
)

var (
    requestCount    = make(map[string]int)
    totalLatency    = make(map[string]time.Duration)
    statusCodeCount = make(map[string]map[int]int)
)

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        path := r.URL.Path

        recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(recorder, r)

        latency := time.Since(start)
        requestCount[path]++
        totalLatency[path] += latency

        if statusCodeCount[path] == nil {
            statusCodeCount[path] = make(map[int]int)
        }
        statusCodeCount[path][recorder.statusCode]++

        fmt.Printf("Request to %s took %v, returned status %d\n", path, latency, recorder.statusCode)
    })
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
    rr.statusCode = code
    rr.ResponseWriter.WriteHeader(code)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintln(w, "HTTP Request Metrics")
    fmt.Fprintln(w, "====================")

    for path, count := range requestCount {
        avgLatency := totalLatency[path] / time.Duration(count)
        fmt.Fprintf(w, "\nPath: %s\n", path)
        fmt.Fprintf(w, "  Total Requests: %d\n", count)
        fmt.Fprintf(w, "  Average Latency: %v\n", avgLatency)
        fmt.Fprintf(w, "  Status Codes:\n")
        for code, freq := range statusCodeCount[path] {
            fmt.Fprintf(w, "    %d: %d\n", code, freq)
        }
    }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(50 * time.Millisecond)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Hello, World!")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(100 * time.Millisecond)
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintln(w, "Internal Server Error")
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/hello", metricsMiddleware(http.HandlerFunc(helloHandler)))
    mux.Handle("/error", metricsMiddleware(http.HandlerFunc(errorHandler)))
    mux.Handle("/metrics", http.HandlerFunc(metricsHandler))

    fmt.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        fmt.Printf("Server failed: %v\n", err)
    }
}