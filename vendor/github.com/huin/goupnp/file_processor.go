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
}package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStats struct {
	Path     string
	LineCount int
	WordCount int
	ByteCount int64
}

func processFile(path string, results chan<- FileStats, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0
	for scanner.Scan() {
		lineCount++
		words := splitWords(scanner.Text())
		wordCount += len(words)
	}

	fileInfo, _ := file.Stat()
	byteCount := fileInfo.Size()

	results <- FileStats{
		Path:      path,
		LineCount: lineCount,
		WordCount: wordCount,
		ByteCount: byteCount,
	}
}

func splitWords(text string) []string {
	var words []string
	wordStart := -1
	for i, r := range text {
		if isWordChar(r) {
			if wordStart == -1 {
				wordStart = i
			}
		} else {
			if wordStart != -1 {
				words = append(words, text[wordStart:i])
				wordStart = -1
			}
		}
	}
	if wordStart != -1 {
		words = append(words, text[wordStart:])
	}
	return words
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> [file2 ...]")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(os.Args)-1)

	for _, filePath := range os.Args[1:] {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Printf("Invalid path %s: %v\n", filePath, err)
			continue
		}

		wg.Add(1)
		go processFile(absPath, results, &wg)
	}

	wg.Wait()
	close(results)

	totalLines := 0
	totalWords := 0
	var totalBytes int64 = 0

	fmt.Println("File Processing Results:")
	fmt.Println("========================")
	for stats := range results {
		fmt.Printf("File: %s\n", stats.Path)
		fmt.Printf("  Lines: %d\n", stats.LineCount)
		fmt.Printf("  Words: %d\n", stats.WordCount)
		fmt.Printf("  Bytes: %d\n", stats.ByteCount)
		fmt.Println()

		totalLines += stats.LineCount
		totalWords += stats.WordCount
		totalBytes += stats.ByteCount
	}

	fmt.Println("Summary:")
	fmt.Printf("Total Files Processed: %d\n", len(os.Args)-1)
	fmt.Printf("Total Lines: %d\n", totalLines)
	fmt.Printf("Total Words: %d\n", totalWords)
	fmt.Printf("Total Bytes: %d\n", totalBytes)
}