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
		CPUUsage:    calculateCPUUsage(),
		MemoryUsage: m.Alloc,
		Goroutines:  runtime.NumGoroutine(),
	}
}

func calculateCPUUsage() float64 {
	start := time.Now()
	runtime.Gosched()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start).Seconds()
	return 1.0 - (elapsed / 0.1)
}

func displayMetrics(metrics SystemMetrics) {
	fmt.Printf("Timestamp: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("CPU Usage: %.2f%%\n", metrics.CPUUsage*100)
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
			displayMetrics(metrics)
		}
	}
}