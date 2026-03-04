package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Processor struct {
	workerCount int
	batchSize   int
}

type DataChunk struct {
	ID   int
	Data string
}

func NewProcessor(workers, batch int) *Processor {
	return &Processor{
		workerCount: workers,
		batchSize:   batch,
	}
}

func (p *Processor) Process(ctx context.Context, data []DataChunk) []string {
	var wg sync.WaitGroup
	results := make([]string, len(data))
	chunkChan := make(chan []DataChunk, p.workerCount)

	for i := 0; i < p.workerCount; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg, chunkChan, results)
	}

	for i := 0; i < len(data); i += p.batchSize {
		end := i + p.batchSize
		if end > len(data) {
			end = len(data)
		}
		chunkChan <- data[i:end]
	}
	close(chunkChan)

	wg.Wait()
	return results
}

func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup, chunks <-chan []DataChunk, results []string) {
	defer wg.Done()

	for chunk := range chunks {
		select {
		case <-ctx.Done():
			return
		default:
			for _, item := range chunk {
				processed := fmt.Sprintf("processed-%d-%s", item.ID, item.Data)
				time.Sleep(10 * time.Millisecond)
				results[item.ID] = processed
			}
		}
	}
}

func main() {
	processor := NewProcessor(4, 10)

	data := make([]DataChunk, 100)
	for i := range data {
		data[i] = DataChunk{ID: i, Data: fmt.Sprintf("item-%d", i)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := processor.Process(ctx, data)

	for i := 0; i < 5; i++ {
		fmt.Println(results[i])
	}
}