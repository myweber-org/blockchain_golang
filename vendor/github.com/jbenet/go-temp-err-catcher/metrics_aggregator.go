package metrics

import (
	"sync"
	"time"
)

type Aggregator struct {
	mu        sync.RWMutex
	window    time.Duration
	metrics   []float64
	timestamps []time.Time
	capacity  int
}

func NewAggregator(window time.Duration, capacity int) *Aggregator {
	return &Aggregator{
		window:    window,
		metrics:   make([]float64, 0, capacity),
		timestamps: make([]time.Time, 0, capacity),
		capacity:  capacity,
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.cleanup(now)

	if len(a.metrics) >= a.capacity {
		a.metrics = a.metrics[1:]
		a.timestamps = a.timestamps[1:]
	}

	a.metrics = append(a.metrics, value)
	a.timestamps = append(a.timestamps, now)
}

func (a *Aggregator) cleanup(now time.Time) {
	cutoff := now.Add(-a.window)
	i := 0
	for i < len(a.timestamps) && a.timestamps[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		a.metrics = a.metrics[i:]
		a.timestamps = a.timestamps[i:]
	}
}

func (a *Aggregator) Average() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	a.cleanup(time.Now())
	if len(a.metrics) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range a.metrics {
		sum += v
	}
	return sum / float64(len(a.metrics))
}

func (a *Aggregator) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	a.cleanup(time.Now())
	return len(a.metrics)
}

func (a *Aggregator) Percentile(p float64) float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	a.cleanup(time.Now())
	if len(a.metrics) == 0 {
		return 0
	}

	values := make([]float64, len(a.metrics))
	copy(values, a.metrics)

	sort.Float64s(values)
	index := int(float64(len(values)-1) * p / 100.0)
	return values[index]
}