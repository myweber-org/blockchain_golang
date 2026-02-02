
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
	baseName    string
	mu          sync.Mutex
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, filename)
	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingWriter{
		currentFile: file,
		currentSize: info.Size(),
		baseName:    filename,
	}, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > maxFileSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.currentFile.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(logDir, w.baseName)
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.%s", w.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.currentFile = file
	w.currentSize = 0

	go w.cleanupOldFiles()
	return nil
}

func (w *RotatingWriter) cleanupOldFiles() {
	pattern := filepath.Join(logDir, w.baseName+".*")
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
	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("application.log")
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(message))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}