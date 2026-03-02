
package main

import (
    "fmt"
    "log"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp   time.Time
    Goroutines  int
    MemoryAlloc uint64
    MemoryTotal uint64
    NumCPU      int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    return SystemMetrics{
        Timestamp:   time.Now(),
        Goroutines:  runtime.NumGoroutine(),
        MemoryAlloc: m.Alloc,
        MemoryTotal: m.TotalAlloc,
        NumCPU:      runtime.NumCPU(),
    }
}

func logMetrics(metrics SystemMetrics) {
    log.Printf(
        "Metrics collected | Time: %s | Goroutines: %d | Alloc: %v MB | TotalAlloc: %v MB | CPUs: %d",
        metrics.Timestamp.Format("2006-01-02 15:04:05"),
        metrics.Goroutines,
        metrics.MemoryAlloc/1024/1024,
        metrics.MemoryTotal/1024/1024,
        metrics.NumCPU,
    )
}

func main() {
    fmt.Println("Starting system metrics collector...")
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            logMetrics(metrics)
        }
    }
}