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
		CPUPercent: getCPUUsage(),
		MemAlloc:   m.Alloc,
		MemTotal:   m.Sys,
		Goroutines: runtime.NumGoroutine(),
	}
}

func getCPUUsage() float64 {
	start := time.Now()
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)
	return float64(elapsed) / float64(time.Second) * 100
}

func printMetrics(metrics SystemMetrics) {
	fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUPercent)
	fmt.Printf("Memory Allocated: %d bytes\n", metrics.MemAlloc)
	fmt.Printf("Total Memory: %d bytes\n", metrics.MemTotal)
	fmt.Printf("Active Goroutines: %d\n", metrics.Goroutines)
	fmt.Println("---")
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 3; i++ {
		metrics := collectMetrics()
		printMetrics(metrics)
		<-ticker.C
	}
}