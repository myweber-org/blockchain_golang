package aggregator

import (
	"container/list"
	"sort"
	"sync"
	"time"
)

type Metric struct {
	Value     float64
	Timestamp time.Time
}

type SlidingWindowAggregator struct {
	windowSize  time.Duration
	maxSamples  int
	metrics     *list.List
	percentiles []float64
	mu          sync.RWMutex
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxSamples int, percentiles []float64) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize:  windowSize,
		maxSamples:  maxSamples,
		metrics:     list.New(),
		percentiles: percentiles,
	}
}

func (swa *SlidingWindowAggregator) AddMetric(value float64) {
	swa.mu.Lock()
	defer swa.mu.Unlock()

	now := time.Now()
	swa.metrics.PushBack(Metric{Value: value, Timestamp: now})

	swa.cleanupOldMetrics(now)
	if swa.metrics.Len() > swa.maxSamples {
		swa.metrics.Remove(swa.metrics.Front())
	}
}

func (swa *SlidingWindowAggregator) cleanupOldMetrics(now time.Time) {
	cutoff := now.Add(-swa.windowSize)
	for e := swa.metrics.Front(); e != nil; {
		next := e.Next()
		if metric := e.Value.(Metric); metric.Timestamp.Before(cutoff) {
			swa.metrics.Remove(e)
		}
		e = next
	}
}

func (swa *SlidingWindowAggregator) GetStats() map[string]float64 {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.metrics.Len() == 0 {
		return make(map[string]float64)
	}

	values := make([]float64, 0, swa.metrics.Len())
	var sum float64
	minVal := 1.7976931348623157e+308
	maxVal := -1.7976931348623157e+308

	for e := swa.metrics.Front(); e != nil; e = e.Next() {
		val := e.Value.(Metric).Value
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
	stats["sum"] = sum
	stats["avg"] = sum / float64(len(values))
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
	return "p" + strings.Replace(fmt.Sprintf("%.1f", p), ".", "_", -1)
}