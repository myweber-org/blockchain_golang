
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingWriter struct {
	mu       sync.Mutex
	filename string
	file     *os.File
	size     int64
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	w := &RotatingWriter{filename: filename}
	if err := w.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *RotatingWriter) rotateIfNeeded() error {
	if w.file == nil || w.size >= maxFileSize {
		if w.file != nil {
			w.file.Close()
		}
		return w.rotate()
	}
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.file != nil {
		oldPath := w.filename
		timestamp := time.Now().Format("20060102-150405")
		newPath := fmt.Sprintf("%s.%s", w.filename, timestamp)
		os.Rename(oldPath, newPath)
		w.cleanupOldFiles()
	}

	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.file = file
	w.size = info.Size()
	return nil
}

func (w *RotatingWriter) cleanupOldFiles() {
	pattern := w.filename + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, file := range toDelete {
			os.Remove(file)
		}
	}
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		writer.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRotator struct {
	filePath     string
	maxSize      int64
	currentSize  int64
	file         *os.File
	mutex        sync.Mutex
	rotationCount int
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	
	return &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		file:        file,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()
	
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log: %w", err)
		}
	}
	
	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}
	
	backupPath := fmt.Sprintf("%s.%d.%s", 
		lr.filePath, 
		lr.rotationCount, 
		time.Now().Format("20060102_150405"))
	
	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}
	
	file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	
	lr.file = file
	lr.currentSize = 0
	lr.rotationCount++
	
	return nil
}

func (lr *LogRotator) Close() error {
	lr.mutex.Lock()
	defer lr.mutex.Unlock()
	
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()
	
	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
			time.Now().Format(time.RFC3339), i)
		
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
		
		time.Sleep(100 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}