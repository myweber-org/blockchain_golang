package main

import (
	"fmt"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Timestamp    time.Time
	MemoryAlloc  uint64
	TotalAlloc   uint64
	Sys          uint64
	NumGC        uint32
	NumGoroutine int
}

func collectMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		Timestamp:    time.Now(),
		MemoryAlloc:  m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		NumGC:        m.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func printMetrics(metrics SystemMetrics) {
	fmt.Printf("Time: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Memory Allocated: %v bytes\n", metrics.MemoryAlloc)
	fmt.Printf("Total Allocated: %v bytes\n", metrics.TotalAlloc)
	fmt.Printf("System Memory: %v bytes\n", metrics.Sys)
	fmt.Printf("Garbage Collections: %v\n", metrics.NumGC)
	fmt.Printf("Goroutines: %v\n", metrics.NumGoroutine)
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