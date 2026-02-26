
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu         sync.Mutex
	current    *os.File
	baseName   string
	fileSize   int64
	sequence   int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logger := &RotatingLogger{
		baseName: baseName,
		sequence: 0,
	}

	if err := logger.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.current.Write(p)
	if err != nil {
		return n, err
	}

	rl.fileSize += int64(n)

	if err := rl.rotateIfNeeded(); err != nil {
		return n, err
	}

	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.current == nil || rl.fileSize >= maxFileSize {
		return rl.rotate()
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.current != nil {
		rl.current.Close()
		if err := rl.compressOldLog(); err != nil {
			fmt.Printf("Failed to compress log: %v\n", err)
		}
	}

	rl.sequence++
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.current = file
	rl.fileSize = 0

	return nil
}

func (rl *RotatingLogger) compressOldLog() error {
	if rl.sequence <= 1 {
		return nil
	}

	oldName := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence-1))
	compressedName := oldName + ".gz"

	oldFile, err := os.Open(oldName)
	if err != nil {
		return err
	}
	defer oldFile.Close()

	compressedFile, err := os.Create(compressedName)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	gzWriter := gzip.NewWriter(compressedFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		return err
	}

	os.Remove(oldName)
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}