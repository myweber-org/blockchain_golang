
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Metrics struct {
	mu            sync.RWMutex
	requestCount  int
	totalLatency  time.Duration
	statusCodes   map[int]int
}

func NewMetrics() *Metrics {
	return &Metrics{
		statusCodes: make(map[int]int),
	}
}

func (m *Metrics) RecordRequest(statusCode int, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount++
	m.totalLatency += latency
	m.statusCodes[statusCode]++
}

func (m *Metrics) GetAverageLatency() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.requestCount == 0 {
		return 0
	}
	return m.totalLatency / time.Duration(m.requestCount)
}

func (m *Metrics) GetStatusCodeDistribution() map[int]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	distribution := make(map[int]int)
	for code, count := range m.statusCodes {
		distribution[code] = count
	}
	return distribution
}

func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount = 0
	m.totalLatency = 0
	m.statusCodes = make(map[int]int)
}

func simulateHTTPRequest(metrics *Metrics) {
	latency := time.Duration(rand.Intn(200)+50) * time.Millisecond
	statusCode := 200
	if rand.Float32() < 0.1 {
		statusCode = 500
	} else if rand.Float32() < 0.05 {
		statusCode = 404
	}

	time.Sleep(latency)
	metrics.RecordRequest(statusCode, latency)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	metrics := NewMetrics()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			simulateHTTPRequest(metrics)
		}()
	}

	wg.Wait()

	fmt.Printf("Total requests: %d\n", metrics.requestCount)
	fmt.Printf("Average latency: %v\n", metrics.GetAverageLatency())
	fmt.Println("Status code distribution:")
	for code, count := range metrics.GetStatusCodeDistribution() {
		fmt.Printf("  %d: %d\n", code, count)
	}
}