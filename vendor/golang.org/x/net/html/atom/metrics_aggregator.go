package main

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

var (
    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    requestCounter = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)

func simulateRequest(method, endpoint string) {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(method, endpoint).Observe(duration)
    }()

    // Simulate random processing time
    processingTime := time.Duration(rand.Intn(200)) * time.Millisecond
    time.Sleep(processingTime)

    // Simulate random status code
    statusCodes := []string{"200", "404", "500"}
    status := statusCodes[rand.Intn(len(statusCodes))]
    requestCounter.WithLabelValues(method, endpoint, status).Inc()

    fmt.Printf("Processed %s %s - Status: %s - Duration: %v\n", method, endpoint, status, processingTime)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    // Start metrics HTTP server
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        fmt.Println("Metrics server listening on :8080")
        http.ListenAndServe(":8080", nil)
    }()

    // Simulate incoming requests
    endpoints := []string{"/api/v1/users", "/api/v1/products", "/api/v1/orders"}
    methods := []string{"GET", "POST", "PUT", "DELETE"}

    for i := 0; i < 100; i++ {
        method := methods[rand.Intn(len(methods))]
        endpoint := endpoints[rand.Intn(len(endpoints))]
        simulateRequest(method, endpoint)
        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    }

    fmt.Println("Request simulation completed. Metrics available at http://localhost:8080/metrics")
    select {} // Keep running
}