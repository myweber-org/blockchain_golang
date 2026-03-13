package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingWriter struct {
	mu          sync.Mutex
	currentSize int64
	maxSize     int64
	basePath    string
	file        *os.File
	sequence    int
}

func NewRotatingWriter(basePath string, maxSize int64) (*RotatingWriter, error) {
	w := &RotatingWriter{
		basePath: basePath,
		maxSize:  maxSize,
		sequence: 0,
	}

	if err := w.openFile(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *RotatingWriter) openFile() error {
	filename := w.basePath
	if w.sequence > 0 {
		filename = fmt.Sprintf("%s.%d", w.basePath, w.sequence)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.file = file
	w.currentSize = stat.Size()
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.file != nil {
		w.file.Close()
	}

	w.sequence++
	return w.openFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
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
	writer, err := NewRotatingWriter("app.log", 1024*1024)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create writer: %v\n", err)
		os.Exit(1)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := writer.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}