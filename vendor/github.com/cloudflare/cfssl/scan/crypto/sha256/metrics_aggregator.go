package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	windowSize time.Duration
	metrics    []float64
	mu         sync.RWMutex
}

func NewAggregator(windowSize time.Duration) *Aggregator {
	return &Aggregator{
		windowSize: windowSize,
		metrics:    make([]float64, 0),
	}
}

func (a *Aggregator) AddMetric(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	a.metrics = append(a.metrics, value)
	a.cleanup()
}

func (a *Aggregator) GetAverage() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.metrics) == 0 {
		return 0.0
	}
	
	var sum float64
	for _, v := range a.metrics {
		sum += v
	}
	return sum / float64(len(a.metrics))
}

func (a *Aggregator) cleanup() {
	cutoff := time.Now().Add(-a.windowSize)
	
	validMetrics := make([]float64, 0)
	for i := len(a.metrics) - 1; i >= 0; i-- {
		if time.Unix(int64(a.metrics[i]), 0).After(cutoff) {
			validMetrics = append([]float64{a.metrics[i]}, validMetrics...)
		}
	}
	
	a.metrics = validMetrics
}

func (a *Aggregator) GetMetricsCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.metrics)
}