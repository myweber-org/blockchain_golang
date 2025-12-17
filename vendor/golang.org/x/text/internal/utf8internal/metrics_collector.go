package main

import (
	"log"
	"net/http"
	"time"
)

type MetricsCollector struct {
	requestCount    int
	errorCount      int
	totalLatency    time.Duration
	latencySamples  []time.Duration
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		latencySamples: make([]time.Duration, 0),
	}
}

func (mc *MetricsCollector) RecordRequest(latency time.Duration, isError bool) {
	mc.requestCount++
	mc.totalLatency += latency
	mc.latencySamples = append(mc.latencySamples, latency)
	
	if isError {
		mc.errorCount++
	}
}

func (mc *MetricsCollector) GetAverageLatency() time.Duration {
	if mc.requestCount == 0 {
		return 0
	}
	return mc.totalLatency / time.Duration(mc.requestCount)
}

func (mc *MetricsCollector) GetErrorRate() float64 {
	if mc.requestCount == 0 {
		return 0.0
	}
	return float64(mc.errorCount) / float64(mc.requestCount)
}

func (mc *MetricsCollector) GetPercentileLatency(percentile float64) time.Duration {
	if len(mc.latencySamples) == 0 {
		return 0
	}
	
	index := int(float64(len(mc.latencySamples)-1) * percentile / 100.0)
	if index < 0 {
		index = 0
	}
	if index >= len(mc.latencySamples) {
		index = len(mc.latencySamples) - 1
	}
	
	return mc.latencySamples[index]
}

func main() {
	collector := NewMetricsCollector()
	
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		
		// Simulate some processing
		time.Sleep(time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond)
		
		latency := time.Since(startTime)
		isError := time.Now().UnixNano()%10 == 0 // 10% error rate simulation
		
		if isError {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Data processed successfully"))
		}
		
		collector.RecordRequest(latency, isError)
	})
	
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics := map[string]interface{}{
			"total_requests": collector.requestCount,
			"error_count":    collector.errorCount,
			"error_rate":     collector.GetErrorRate(),
			"avg_latency_ms": collector.GetAverageLatency().Milliseconds(),
			"p95_latency_ms": collector.GetPercentileLatency(95).Milliseconds(),
			"p99_latency_ms": collector.GetPercentileLatency(99).Milliseconds(),
		}
		
		jsonData := `{"total_requests":%d,"error_count":%d,"error_rate":%.4f,"avg_latency_ms":%.2f,"p95_latency_ms":%.2f,"p99_latency_ms":%.2f}`
		response := fmt.Sprintf(jsonData,
			metrics["total_requests"],
			metrics["error_count"],
			metrics["error_rate"],
			metrics["avg_latency_ms"],
			metrics["p95_latency_ms"],
			metrics["p99_latency_ms"])
		
		w.Write([]byte(response))
	})
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}