package main

import (
    "fmt"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp     time.Time
    CPUUsage      float64
    MemoryUsageMB float64
    GoroutineCount int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:     time.Now(),
        MemoryUsageMB: float64(m.Alloc) / 1024 / 1024,
        GoroutineCount: runtime.NumGoroutine(),
    }
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        metrics := collectMetrics()
        fmt.Printf("[%s] Memory: %.2f MB, Goroutines: %d\n",
            metrics.Timestamp.Format("15:04:05"),
            metrics.MemoryUsageMB,
            metrics.GoroutineCount)
    }
}