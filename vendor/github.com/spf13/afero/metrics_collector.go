package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "runtime"
    "time"
)

type SystemMetrics struct {
    Timestamp    time.Time `json:"timestamp"`
    MemoryAlloc  uint64    `json:"memory_alloc_bytes"`
    TotalAlloc   uint64    `json:"total_alloc_bytes"`
    SysAlloc     uint64    `json:"sys_alloc_bytes"`
    NumGoroutine int       `json:"goroutine_count"`
    NumCPU       int       `json:"cpu_count"`
}

func collectMetrics() SystemMetrics {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    return SystemMetrics{
        Timestamp:    time.Now().UTC(),
        MemoryAlloc:  memStats.Alloc,
        TotalAlloc:   memStats.TotalAlloc,
        SysAlloc:     memStats.Sys,
        NumGoroutine: runtime.NumGoroutine(),
        NumCPU:       runtime.NumCPU(),
    }
}

func writeMetricsToFile(metrics SystemMetrics, filename string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(metrics)
}

func main() {
    const outputFile = "system_metrics.jsonl"
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    log.Printf("Starting metrics collection. Output file: %s\n", outputFile)

    for {
        select {
        case <-ticker.C:
            metrics := collectMetrics()
            if err := writeMetricsToFile(metrics, outputFile); err != nil {
                log.Printf("Failed to write metrics: %v\n", err)
                continue
            }
            fmt.Printf("Metrics collected at %s: MemoryAlloc=%d bytes, Goroutines=%d\n",
                metrics.Timestamp.Format(time.RFC3339),
                metrics.MemoryAlloc,
                metrics.NumGoroutine)
        }
    }
}