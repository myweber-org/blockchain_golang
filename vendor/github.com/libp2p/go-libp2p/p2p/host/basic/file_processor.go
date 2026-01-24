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

type FileStats struct {
	Path        string
	Size        int64
	LineCount   int
	ProcessTime time.Duration
}

func processFile(path string, results chan<- FileStats, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	stats := FileStats{Path: path}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info for %s: %v\n", path, err)
		return
	}
	stats.Size = fileInfo.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return
	}
	stats.LineCount = lineCount

	stats.ProcessTime = time.Since(start)
	results <- stats
}

func collectFiles(dir string, patterns []string) ([]string, error) {
	var files []string
	visited := make(map[string]bool)

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			if !visited[match] {
				info, err := os.Stat(match)
				if err != nil {
					continue
				}
				if !info.IsDir() {
					files = append(files, match)
					visited[match] = true
				}
			}
		}
	}

	if len(files) == 0 {
		return nil, errors.New("no matching files found")
	}
	return files, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory> [file_patterns...]")
		fmt.Println("Default patterns: *.txt, *.go, *.md")
		os.Exit(1)
	}

	dir := os.Args[1]
	patterns := []string{"*.txt", "*.go", "*.md"}
	if len(os.Args) > 2 {
		patterns = os.Args[2:]
	}

	files, err := collectFiles(dir, patterns)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing %d files from directory: %s\n", len(files), dir)

	results := make(chan FileStats, len(files))
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go processFile(file, results, &wg)
	}

	wg.Wait()
	close(results)

	var totalSize int64
	var totalLines int
	var totalTime time.Duration

	fmt.Println("\nFile Processing Results:")
	fmt.Println("========================")
	for stats := range results {
		fmt.Printf("File: %s\n", filepath.Base(stats.Path))
		fmt.Printf("  Size: %d bytes\n", stats.Size)
		fmt.Printf("  Lines: %d\n", stats.LineCount)
		fmt.Printf("  Process Time: %v\n", stats.ProcessTime)
		fmt.Println()

		totalSize += stats.Size
		totalLines += stats.LineCount
		totalTime += stats.ProcessTime
	}

	fmt.Println("Summary:")
	fmt.Println("========")
	fmt.Printf("Total Files Processed: %d\n", len(files))
	fmt.Printf("Total Size: %d bytes\n", totalSize)
	fmt.Printf("Total Lines: %d\n", totalLines)
	fmt.Printf("Total Processing Time: %v\n", totalTime)
	if len(files) > 0 {
		fmt.Printf("Average Time per File: %v\n", totalTime/time.Duration(len(files)))
	}
}