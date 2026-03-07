package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUUsage    float64
    MemoryUsage uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now(),
        CPUUsage:    getCPUUsage(),
        MemoryUsage: m.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func getCPUUsage() float64 {
    // Simplified CPU usage calculation
    // In production, use proper OS-specific metrics collection
    return float64(runtime.NumCPU()) * 0.75 // Simulated 75% usage
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
    fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
    fmt.Printf("Memory Usage: %d bytes\n", metrics.MemoryUsage)
    fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
    fmt.Println("---")
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}