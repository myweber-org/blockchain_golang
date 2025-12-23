package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemoryAlloc uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:   time.Now(),
        MemoryAlloc: m.Alloc,
        Goroutines:  runtime.NumGoroutine(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("Time: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
    fmt.Printf("Memory Usage: %d bytes\n", metrics.MemoryAlloc)
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