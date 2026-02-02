package main

import (
	"bufio"
	"fmt"
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

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Printf("Error stating %s: %v\n", path, err)
		return
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning %s: %v\n", path, err)
		return
	}

	results <- FileStats{
		Path:     path,
		Size:     stat.Size(),
		Lines:    lineCount,
		Modified: stat.ModTime(),
	}
}

func collectFiles(dir string, extensions []string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			for _, ext := range extensions {
				if filepath.Ext(path) == ext {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})

	return files, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	extensions := []string{".txt", ".go", ".md", ".json"}

	files, err := collectFiles(dir, extensions)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(files))

	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSize := int64(0)
	totalLines := 0
	fileCount := 0

	for stats := range results {
		fmt.Printf("File: %s\n", stats.Path)
		fmt.Printf("  Size: %d bytes\n", stats.Size)
		fmt.Printf("  Lines: %d\n", stats.Lines)
		fmt.Printf("  Modified: %s\n\n", stats.Modified.Format("2006-01-02 15:04:05"))

		totalSize += stats.Size
		totalLines += stats.Lines
		fileCount++
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("  Files processed: %d\n", fileCount)
	fmt.Printf("  Total size: %d bytes\n", totalSize)
	fmt.Printf("  Total lines: %d\n", totalLines)
}