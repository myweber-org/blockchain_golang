
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