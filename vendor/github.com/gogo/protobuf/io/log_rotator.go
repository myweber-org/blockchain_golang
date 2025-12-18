package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	fileCount   int
	maxFiles    int
	currentSize int64
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	rl.fileCount = 1

	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile == nil {
		return fmt.Errorf("no current file open")
	}

	rl.currentFile.Close()

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

	if err := os.Rename(rl.basePath, rotatedPath); err != nil {
		return err
	}

	if err := rl.compressFile(rotatedPath); err != nil {
		fmt.Printf("Warning: failed to compress %s: %v\n", rotatedPath, err)
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	rl.fileCount++

	if rl.fileCount > rl.maxFiles {
		rl.cleanupOldFiles()
	}

	return nil
}

func (rl *RotatingLogger) compressFile(sourcePath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destPath := sourcePath + ".gz"
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, source); err != nil {
		return err
	}

	os.Remove(sourcePath)
	return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxFiles {
		filesToRemove := len(matches) - rl.maxFiles
		for i := 0; i < filesToRemove; i++ {
			os.Remove(matches[i])
		}
	}
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		rl.mu.Unlock()
		if err := rl.rotate(); err != nil {
			rl.mu.Lock()
			return 0, err
		}
		rl.mu.Lock()
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}