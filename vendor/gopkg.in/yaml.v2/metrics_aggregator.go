
package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize time.Duration
	percentile float64
	mu         sync.RWMutex
	values     []float64
	timestamps []time.Time
}

func NewAggregator(windowSize time.Duration, percentile float64) *Aggregator {
	return &Aggregator{
		windowSize: windowSize,
		percentile: percentile,
		values:     make([]float64, 0),
		timestamps: make([]time.Time, 0),
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	now := time.Now()
	a.values = append(a.values, value)
	a.timestamps = append(a.timestamps, now)
	a.cleanup(now)
}

func (a *Aggregator) GetPercentile() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	if len(a.values) == 0 {
		return 0.0
	}
	
	values := make([]float64, len(a.values))
	copy(values, a.values)
	sort.Float64s(values)
	
	index := int(a.percentile * float64(len(values)-1))
	return values[index]
}

func (a *Aggregator) cleanup(currentTime time.Time) {
	cutoff := currentTime.Add(-a.windowSize)
	i := 0
	for ; i < len(a.timestamps); i++ {
		if a.timestamps[i].After(cutoff) {
			break
		}
	}
	
	if i > 0 {
		a.values = a.values[i:]
		a.timestamps = a.timestamps[i:]
	}
}