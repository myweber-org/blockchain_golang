
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
}

func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		mc.requestCount++
		mc.totalLatency += duration
		
		if recorder.statusCode >= 400 {
			mc.errorCount++
		}
		
		log.Printf("Request processed: %s %s - Status: %d, Duration: %v", 
			r.Method, r.URL.Path, recorder.statusCode, duration)
	})
}

func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
	avgLatency := time.Duration(0)
	if mc.requestCount > 0 {
		avgLatency = mc.totalLatency / time.Duration(mc.requestCount)
	}
	
	errorRate := 0.0
	if mc.requestCount > 0 {
		errorRate = float64(mc.errorCount) / float64(mc.requestCount) * 100
	}
	
	return map[string]interface{}{
		"total_requests":   mc.requestCount,
		"error_count":      mc.errorCount,
		"error_rate":       errorRate,
		"avg_latency_ms":   avgLatency.Milliseconds(),
		"total_latency_ms": mc.totalLatency.Milliseconds(),
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func main() {
	collector := &MetricsCollector{}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := collector.GetMetrics()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	})
	
	handler := collector.Middleware(mux)
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}