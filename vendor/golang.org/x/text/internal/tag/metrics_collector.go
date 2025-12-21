package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp   time.Time
	CPUUsage    float64
	MemoryAlloc uint64
	MemoryTotal uint64
	Goroutines  int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:   time.Now(),
		CPUUsage:    getCPUUsage(),
		MemoryAlloc: m.Alloc,
		MemoryTotal: m.Sys,
		Goroutines:  runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(50 * time.Millisecond)
	elapsed := time.Since(start).Seconds()
	return (50.0 / 1000.0) / elapsed * 100.0
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format(time.RFC3339))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage)
	fmt.Printf("Memory Allocated: %v bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Memory: %v bytes\n", metrics.MemoryTotal)
	fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
	fmt.Println("---")
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		<-ticker.C
		metrics := collectMetrics()
		displayMetrics(metrics)
	}
}