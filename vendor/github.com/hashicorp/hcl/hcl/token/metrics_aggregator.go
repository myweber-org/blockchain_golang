
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
}

func NewSlidingWindowAggregator(windowSize time.Duration, maxPoints int) *SlidingWindowAggregator {
	return &SlidingWindowAggregator{
		windowSize: windowSize,
		maxPoints:  maxPoints,
		points:     list.New(),
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

	swa.cleanup(now)
}

func (swa *SlidingWindowAggregator) cleanup(currentTime time.Time) {
	cutoff := currentTime.Add(-swa.windowSize)

	for swa.points.Len() > 0 {
		front := swa.points.Front()
		if front == nil {
			break
		}

		point := front.Value.(MetricPoint)
		if point.Timestamp.After(cutoff) && swa.points.Len() <= swa.maxPoints {
			break
		}

		swa.points.Remove(front)
	}
}

func (swa *SlidingWindowAggregator) GetPercentile(p float64) (float64, bool) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.points.Len() == 0 {
		return 0, false
	}

	values := make([]float64, 0, swa.points.Len())
	for e := swa.points.Front(); e != nil; e = e.Next() {
		values = append(values, e.Value.(MetricPoint).Value)
	}

	sort.Float64s(values)

	index := p * float64(len(values)-1) / 100.0
	lower := int(index)
	upper := lower + 1

	if upper >= len(values) {
		return values[lower], true
	}

	weight := index - float64(lower)
	return values[lower]*(1-weight) + values[upper]*weight, true
}

func (swa *SlidingWindowAggregator) GetStats() (min, max, avg float64, count int) {
	swa.mu.RLock()
	defer swa.mu.RUnlock()

	if swa.points.Len() == 0 {
		return 0, 0, 0, 0
	}

	var sum float64
	min = swa.points.Front().Value.(MetricPoint).Value
	max = min

	for e := swa.points.Front(); e != nil; e = e.Next() {
		value := e.Value.(MetricPoint).Value
		sum += value
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
		count++
	}

	avg = sum / float64(count)
	return min, max, avg, count
}