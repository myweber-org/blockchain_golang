
package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	mu               sync.RWMutex
	requestCount     int
	totalLatency     time.Duration
	statusCodeCounts map[int]int
}

func NewMetrics() *Metrics {
	return &Metrics{
		statusCodeCounts: make(map[int]int),
	}
}

func (m *Metrics) RecordRequest(statusCode int, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCount++
	m.totalLatency += latency
	m.statusCodeCounts[statusCode]++
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
	dist := make(map[int]int)
	for k, v := range m.statusCodeCounts {
		dist[k] = v
	}
	return dist
}

func main() {
	metrics := NewMetrics()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		statusCode := http.StatusOK

		defer func() {
			latency := time.Since(start)
			metrics.RecordRequest(statusCode, latency)
		}()

		w.WriteHeader(statusCode)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		avgLatency := metrics.GetAverageLatency()
		distribution := metrics.GetStatusCodeDistribution()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"average_latency_ms": avgLatency.Milliseconds(),
			"status_codes":       distribution,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}