
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Metrics struct {
	mu            sync.RWMutex
	requestCount  int
	totalLatency  time.Duration
	statusCodes   map[int]int
}

func NewMetrics() *Metrics {
	return &Metrics{
		statusCodes: make(map[int]int),
	}
}

func (m *Metrics) RecordRequest(statusCode int, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount++
	m.totalLatency += latency
	m.statusCodes[statusCode]++
}

func (m *Metrics) GetAverageLatency() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.requestCount == 0 {
		return 0
	}
	return m.totalLatency / time.Duration(m.requestCount)
}

func (m *Metrics) GetStatusCodeDistribution() map[int]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	distribution := make(map[int]int)
	for code, count := range m.statusCodes {
		distribution[code] = count
	}
	return distribution
}

func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount = 0
	m.totalLatency = 0
	m.statusCodes = make(map[int]int)
}

func simulateHTTPRequest(metrics *Metrics) {
	latency := time.Duration(rand.Intn(200)+50) * time.Millisecond
	statusCode := 200
	if rand.Float32() < 0.1 {
		statusCode = 500
	} else if rand.Float32() < 0.05 {
		statusCode = 404
	}

	time.Sleep(latency)
	metrics.RecordRequest(statusCode, latency)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	metrics := NewMetrics()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			simulateHTTPRequest(metrics)
		}()
	}

	wg.Wait()

	fmt.Printf("Total requests: %d\n", metrics.requestCount)
	fmt.Printf("Average latency: %v\n", metrics.GetAverageLatency())
	fmt.Println("Status code distribution:")
	for code, count := range metrics.GetStatusCodeDistribution() {
		fmt.Printf("  %d: %d\n", code, count)
	}
}package metrics

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	windowSize   time.Duration
	maxSamples   int
	measurements []measurement
	mu           sync.RWMutex
}

type measurement struct {
	timestamp time.Time
	value     float64
}

type Summary struct {
	Count    int
	Mean     float64
	Median   float64
	P95      float64
	P99      float64
	Min      float64
	Max      float64
}

func NewAggregator(windowSize time.Duration, maxSamples int) *Aggregator {
	return &Aggregator{
		windowSize: windowSize,
		maxSamples: maxSamples,
		measurements: make([]measurement, 0, maxSamples),
	}
}

func (a *Aggregator) Record(value float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	a.measurements = append(a.measurements, measurement{
		timestamp: now,
		value:     value,
	})

	a.prune(now)
	if len(a.measurements) > a.maxSamples {
		a.measurements = a.measurements[1:]
	}
}

func (a *Aggregator) prune(now time.Time) {
	cutoff := now.Add(-a.windowSize)
	i := 0
	for i < len(a.measurements) && a.measurements[i].timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		a.measurements = a.measurements[i:]
	}
}

func (a *Aggregator) GetSummary() Summary {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.measurements) == 0 {
		return Summary{}
	}

	values := make([]float64, len(a.measurements))
	var sum float64
	min := a.measurements[0].value
	max := a.measurements[0].value

	for i, m := range a.measurements {
		values[i] = m.value
		sum += m.value
		if m.value < min {
			min = m.value
		}
		if m.value > max {
			max = m.value
		}
	}

	sort.Float64s(values)

	return Summary{
		Count:  len(values),
		Mean:   sum / float64(len(values)),
		Median: percentile(values, 0.5),
		P95:    percentile(values, 0.95),
		P99:    percentile(values, 0.99),
		Min:    min,
		Max:    max,
	}
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	index := p * float64(len(sorted)-1)
	lower := int(index)
	upper := lower + 1
	weight := index - float64(lower)

	if upper >= len(sorted) {
		return sorted[lower]
	}
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}