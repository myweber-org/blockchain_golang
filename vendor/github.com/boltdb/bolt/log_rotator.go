
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
	logDir      = "./logs"
)

type RotatingWriter struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingWriter(name string) (*RotatingWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	w := &RotatingWriter{baseName: basePath}

	if err := w.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := w.currentFile.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) rotateIfNeeded() error {
	if w.currentFile == nil || w.currentSize >= maxFileSize {
		return w.rotate()
	}
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.currentFile != nil {
		if err := w.currentFile.Close(); err != nil {
			return err
		}
		w.compressOldFile(w.currentFile.Name())
	}

	timestamp := time.Now().Format("20060102_150405")
	newPath := fmt.Sprintf("%s_%s.log", w.baseName, timestamp)

	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.currentFile = file
	if stat, err := file.Stat(); err == nil {
		w.currentSize = stat.Size()
	} else {
		w.currentSize = 0
	}

	w.cleanupOldFiles()
	return nil
}

func (w *RotatingWriter) compressOldFile(path string) {
	// Compression logic placeholder
	// In production, implement actual compression here
	fmt.Printf("Would compress: %s\n", path)
}

func (w *RotatingWriter) cleanupOldFiles() {
	pattern := fmt.Sprintf("%s_*.log", w.baseName)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		filesToRemove := matches[:len(matches)-maxBackups]
		for _, file := range filesToRemove {
			os.Remove(file)
		}
	}
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app")
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(logEntry))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
}