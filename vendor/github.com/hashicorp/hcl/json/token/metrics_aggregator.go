
package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize   time.Duration
	dataPoints   []float64
	timestamps   []time.Time
	mu           sync.RWMutex
	percentiles  []float64
}

func NewAggregator(windowSize time.Duration, percentiles []float64) *Aggregator {
	return &Aggregator{
		windowSize:  windowSize,
		percentiles: percentiles,
		dataPoints:  make([]float64, 0),
		timestamps:  make([]time.Time, 0),
	}
}

func (a *Aggregator) Add(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	now := time.Now()
	a.dataPoints = append(a.dataPoints, value)
	a.timestamps = append(a.timestamps, now)
	
	a.cleanup(now)
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
		a.dataPoints = a.dataPoints[i:]
		a.timestamps = a.timestamps[i:]
	}
}

func (a *Aggregator) GetStats() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	a.cleanup(time.Now())
	
	if len(a.dataPoints) == 0 {
		return make(map[string]float64)
	}
	
	stats := make(map[string]float64)
	
	sortedPoints := make([]float64, len(a.dataPoints))
	copy(sortedPoints, a.dataPoints)
	sort.Float64s(sortedPoints)
	
	var sum float64
	for _, v := range a.dataPoints {
		sum += v
	}
	
	stats["count"] = float64(len(a.dataPoints))
	stats["mean"] = sum / stats["count"]
	stats["min"] = sortedPoints[0]
	stats["max"] = sortedPoints[len(sortedPoints)-1]
	
	for _, p := range a.percentiles {
		key := formatPercentileKey(p)
		stats[key] = calculatePercentile(sortedPoints, p)
	}
	
	return stats
}

func calculatePercentile(sortedData []float64, percentile float64) float64 {
	if len(sortedData) == 0 {
		return 0
	}
	
	index := (percentile / 100) * float64(len(sortedData)-1)
	
	lower := int(index)
	upper := lower + 1
	
	if upper >= len(sortedData) {
		return sortedData[lower]
	}
	
	weight := index - float64(lower)
	return sortedData[lower]*(1-weight) + sortedData[upper]*weight
}

func formatPercentileKey(p float64) string {
	return "p" + formatFloat(p)
}

func formatFloat(f float64) string {
	s := formatFloatNoTrailingZeros(f)
	if s == "0" {
		return "0"
	}
	return s
}

func formatFloatNoTrailingZeros(f float64) string {
	s := formatFloatWithPrecision(f, 2)
	for s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}

func formatFloatWithPrecision(f float64, precision int) string {
	format := "%." + formatInt(precision) + "f"
	return fmt.Sprintf(format, f)
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}