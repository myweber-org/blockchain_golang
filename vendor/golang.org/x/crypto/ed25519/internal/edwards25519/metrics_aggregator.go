package metrics

import (
	"sync"
	"time"
)

type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

type Metric struct {
	Name  string
	Type  MetricType
	Value float64
	Tags  map[string]string
	Time  time.Time
}

type SlidingWindow struct {
	windowSize time.Duration
	metrics    []Metric
	mu         sync.RWMutex
}

func NewSlidingWindow(windowSize time.Duration) *SlidingWindow {
	return &SlidingWindow{
		windowSize: windowSize,
		metrics:    make([]Metric, 0),
	}
}

func (sw *SlidingWindow) AddMetric(metric Metric) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	
	metric.Time = time.Now()
	sw.metrics = append(sw.metrics, metric)
	sw.cleanup()
}

func (sw *SlidingWindow) cleanup() {
	cutoff := time.Now().Add(-sw.windowSize)
	validStart := 0
	
	for i, metric := range sw.metrics {
		if metric.Time.After(cutoff) {
			validStart = i
			break
		}
	}
	
	if validStart > 0 {
		sw.metrics = sw.metrics[validStart:]
	}
}

func (sw *SlidingWindow) Aggregate(metricName string, operation string) float64 {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	
	sw.cleanup()
	
	var result float64
	count := 0
	
	for _, metric := range sw.metrics {
		if metric.Name != metricName {
			continue
		}
		
		switch operation {
		case "sum":
			result += metric.Value
		case "avg":
			result += metric.Value
			count++
		case "max":
			if metric.Value > result {
				result = metric.Value
			}
		case "min":
			if count == 0 || metric.Value < result {
				result = metric.Value
			}
		}
		count++
	}
	
	if operation == "avg" && count > 0 {
		return result / float64(count)
	}
	
	return result
}

func (sw *SlidingWindow) GetMetrics() []Metric {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	
	sw.cleanup()
	return append([]Metric{}, sw.metrics...)
}

func (sw *SlidingWindow) Clear() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.metrics = make([]Metric, 0)
}