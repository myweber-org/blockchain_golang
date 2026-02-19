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
	Path     string
	Size     int64
	Lines    int
	Modified time.Time
}

func processFile(path string) (*FileStats, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &FileStats{
		Path:     path,
		Size:     stat.Size(),
		Lines:    lineCount,
		Modified: stat.ModTime(),
	}, nil
}

func processFilesConcurrently(paths []string) ([]FileStats, []error) {
	var wg sync.WaitGroup
	statsChan := make(chan FileStats, len(paths))
	errChan := make(chan error, len(paths))

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			stats, err := processFile(p)
			if err != nil {
				errChan <- fmt.Errorf("processing %s: %w", p, err)
				return
			}
			statsChan <- *stats
		}(path)
	}

	wg.Wait()
	close(statsChan)
	close(errChan)

	var results []FileStats
	for stat := range statsChan {
		results = append(results, stat)
	}

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	return results, errors
}

func findFilesByExtension(dir, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func writeStatsToFile(stats []FileStats, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, stat := range stats {
		line := fmt.Sprintf("%s|%d|%d|%s\n",
			stat.Path,
			stat.Size,
			stat.Lines,
			stat.Modified.Format(time.RFC3339))
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: file_processor <directory> <extension>")
		os.Exit(1)
	}

	dir := os.Args[1]
	ext := os.Args[2]

	files, err := findFilesByExtension(dir, ext)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found with the specified extension")
		return
	}

	fmt.Printf("Processing %d files...\n", len(files))
	stats, errors := processFilesConcurrently(files)

	if len(errors) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	if len(stats) > 0 {
		outputFile := "file_stats.txt"
		if err := writeStatsToFile(stats, outputFile); err != nil {
			fmt.Printf("Error writing output: %v\n", err)
		} else {
			fmt.Printf("Results written to %s\n", outputFile)
		}

		totalSize := int64(0)
		totalLines := 0
		for _, stat := range stats {
			totalSize += stat.Size
			totalLines += stat.Lines
		}
		fmt.Printf("Total: %d files, %d bytes, %d lines\n",
			len(stats), totalSize, totalLines)
	}
}