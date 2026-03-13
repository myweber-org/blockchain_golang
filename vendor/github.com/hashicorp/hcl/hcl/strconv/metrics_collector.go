package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    CPUPercent  float64
    MemAlloc    uint64
    MemTotal    uint64
    Goroutines  int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:  time.Now(),
        MemAlloc:   m.Alloc,
        MemTotal:   m.Sys,
        Goroutines: runtime.NumGoroutine(),
    }
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("[%s] Goroutines: %d, Memory: %v/%v bytes\n",
        metrics.Timestamp.Format("15:04:05"),
        metrics.Goroutines,
        metrics.MemAlloc,
        metrics.MemTotal)
}

func main() {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}