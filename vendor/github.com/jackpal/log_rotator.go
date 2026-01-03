
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
	backupCount = 5
)

type RotatingWriter struct {
	filename   string
	current    *os.File
	size       int64
	mu         sync.Mutex
	maxSize    int64
	maxBackups int
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	w := &RotatingWriter{
		filename:   filename,
		maxSize:    maxFileSize,
		maxBackups: backupCount,
	}

	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openFile() error {
	dir := filepath.Dir(w.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.current = file
	w.size = stat.Size()
	return nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) >= w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.current.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if w.current != nil {
		if err := w.current.Close(); err != nil {
			return err
		}
	}

	for i := w.maxBackups - 1; i >= 0; i-- {
		oldName := w.backupName(i)
		newName := w.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(w.filename, w.backupName(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return w.openFile()
}

func (w *RotatingWriter) backupName(i int) string {
	if i == 0 {
		return w.filename + ".1"
	}
	return fmt.Sprintf("%s.%d", w.filename, i+1)
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.current != nil {
		return w.current.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log")
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := writer.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write failed: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}