
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers    int
	bufferSize int
	mu         sync.RWMutex
	stats      map[string]int
}

func NewFileProcessor(workers, bufferSize int) *FileProcessor {
	return &FileProcessor{
		workers:    workers,
		bufferSize: bufferSize,
		stats:      make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(path string) error {
	if !fp.isValidPath(path) {
		return errors.New("invalid file path")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, fp.bufferSize)
	scanner.Buffer(buf, fp.bufferSize)

	lines := make(chan string, fp.workers*2)
	results := make(chan int, fp.workers*2)
	var wg sync.WaitGroup

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(lines, results, &wg)
	}

	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for count := range results {
		total += count
	}

	fp.mu.Lock()
	fp.stats[filepath.Base(path)] = total
	fp.mu.Unlock()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func (fp *FileProcessor) worker(lines <-chan string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines {
		results <- fp.processLine(line)
	}
}

func (fp *FileProcessor) processLine(line string) int {
	time.Sleep(1 * time.Millisecond)
	return len(line)
}

func (fp *FileProcessor) isValidPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (fp *FileProcessor) GetStats() map[string]int {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	statsCopy := make(map[string]int, len(fp.stats))
	for k, v := range fp.stats {
		statsCopy[k] = v
	}
	return statsCopy
}

func main() {
	processor := NewFileProcessor(4, 64*1024)

	files := []string{"test1.txt", "test2.txt"}
	for _, file := range files {
		if err := processor.ProcessFile(file); err != nil {
			fmt.Printf("Error processing %s: %v\n", file, err)
		}
	}

	fmt.Println("Processing statistics:")
	for file, count := range processor.GetStats() {
		fmt.Printf("%s: %d characters processed\n", file, count)
	}
}