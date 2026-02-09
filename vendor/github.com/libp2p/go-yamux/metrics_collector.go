
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
    start := time.Now()
    runtime.Gosched()
    elapsed := time.Since(start).Seconds()
    return elapsed * 100
}

func monitorSystem(interval time.Duration, stopChan chan bool) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            fmt.Printf("[%s] CPU: %.2f%%, Memory: %v bytes, Goroutines: %d\n",
                metrics.Timestamp.Format("15:04:05"),
                metrics.CPUUsage,
                metrics.MemoryUsage,
                metrics.Goroutines)
        case <-stopChan:
            fmt.Println("Monitoring stopped")
            return
        }
    }
}

func main() {
    stopChan := make(chan bool)
    go monitorSystem(2*time.Second, stopChan)

    time.Sleep(10 * time.Second)
    stopChan <- true
    time.Sleep(500 * time.Millisecond)
}