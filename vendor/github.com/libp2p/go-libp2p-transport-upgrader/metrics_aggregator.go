
package metrics

import (
	"container/list"
	"sort"
	"sync"
	"time"
)

type MetricPoint struct {
	Value     float64
	Timestamp time.Time
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	maxPoints   int
	points      *list.List
	mu          sync.RWMutex
	percentiles []float64
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxPoints int, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		maxPoints:   maxPoints,
		points:      list.New(),
		percentiles: percentiles,
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.points.PushBack(MetricPoint{
		Value:     value,
		Timestamp: now,
	})

	swa.cleanupOldPoints(now)
	if swa.points.Len() > swa.maxPoints {
		swa.points.Remove(swa.points.Front())
	}
}

func (swa *SlidingWindowAggregator) cleanupOldPoints(now time.Time) {
	cutoff := now.Add(-swa.windowSize)
	for e := swa.points.Front(); e != nil; {
		next := e.Next()
		if mp := e.Value.(MetricPoint); mp.Timestamp.Before(cutoff) {
			swa.points.Remove(e)
		}
		e = next
	}
}

func (swa *SlidingWindowAggregator) GetStats() map[string]float64 {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.points.Len() == 0 {
		return make(map[string]float64)
	}

	values := make([]float64, 0, swa.points.Len())
	var sum float64
	minVal := 1e100
	maxVal := -1e100

	for e := swa.points.Front(); e != nil; e = e.Next() {
		val := e.Value.(MetricPoint).Value
		values = append(values, val)
		sum += val
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	sort.Float64s(values)
	stats := make(map[string]float64)
	stats["count"] = float64(len(values))
	stats["mean"] = sum / float64(len(values))
	stats["min"] = minVal
	stats["max"] = maxVal

	for _, p := range swa.percentiles {
		if p < 0 || p > 100 {
			continue
		}
		idx := int(float64(len(values)-1) * p / 100.0)
		stats[formatPercentileKey(p)] = values[idx]
	}

	return stats
}

func formatPercentileKey(p float64) string {
	if p == 50 {
		return "median"
	}
	return string(rune('p')) + formatFloat(p)
}

func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return string(rune(int(f)))
	}
	return string(rune(f))
}