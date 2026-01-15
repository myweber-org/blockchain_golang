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
	rotation    int
}

func NewRotatingWriter(basePath string, maxSize int64) (*RotatingWriter, error) {
	w := &RotatingWriter{
		maxSize:  maxSize,
		basePath: basePath,
	}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openFile() error {
	path := w.basePath
	if w.rotation > 0 {
		path = fmt.Sprintf("%s.%d", w.basePath, w.rotation)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
	w.file.Close()
	w.rotation++
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
	return w.file.Close()
}

func main() {
	writer, err := NewRotatingWriter("app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}