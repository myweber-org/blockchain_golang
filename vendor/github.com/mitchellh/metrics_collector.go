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
    NumGoroutine int
}

func collectMetrics() SystemMetrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return SystemMetrics{
        Timestamp:   time.Now().UTC(),
        CPUPercent:  getCPUUsage(),
        MemoryAlloc: m.Alloc,
        NumGoroutine: runtime.NumGoroutine(),
    }
}

func getCPUUsage() float64 {
    start := time.Now()
    startCPU := runtime.NumCPU()
    
    time.Sleep(100 * time.Millisecond)
    
    elapsed := time.Since(start)
    endCPU := runtime.NumCPU()
    
    if elapsed == 0 || startCPU == 0 {
        return 0.0
    }
    
    usage := float64(endCPU-startCPU) / float64(startCPU) * 100
    if usage < 0 {
        usage = 0
    }
    return usage
}

func printMetrics(metrics SystemMetrics) {
    fmt.Printf("[%s] CPU: %.2f%% | Memory: %v bytes | Goroutines: %d\n",
        metrics.Timestamp.Format("2006-01-02 15:04:05"),
        metrics.CPUPercent,
        metrics.MemoryAlloc,
        metrics.NumGoroutine)
}

func main() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    fmt.Println("Starting system metrics collector...")
    
    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            printMetrics(metrics)
        }
    }
}