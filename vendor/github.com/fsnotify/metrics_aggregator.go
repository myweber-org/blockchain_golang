package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	windowSize time.Duration
	mu         sync.RWMutex
	metrics    []float64
	timestamps []time.Time
}

func NewAggregator(windowSize time.Duration) *Aggregator {
	return &Aggregator{
		windowSize: windowSize,
		metrics:    make([]float64, 0),
		timestamps: make([]time.Time, 0),
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	now := time.Now()
	a.metrics = append(a.metrics, value)
	a.timestamps = append(a.timestamps, now)
	a.cleanup(now)
}

func (a *Aggregator) GetAverage() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.metrics) == 0 {
		return 0
	}
	
	a.cleanup(time.Now())
	
	var sum float64
	for _, v := range a.metrics {
		sum += v
	}
	return sum / float64(len(a.metrics))
}

func (a *Aggregator) cleanup(currentTime time.Time) {
	cutoff := currentTime.Add(-a.windowSize)
	
	i := 0
	for i < len(a.timestamps) && a.timestamps[i].Before(cutoff) {
		i++
	}
	
	if i > 0 {
		a.metrics = a.metrics[i:]
		a.timestamps = a.timestamps[i:]
	}
}