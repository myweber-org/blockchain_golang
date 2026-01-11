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
	maxFileSize = 1024 * 1024
	maxBackups  = 5
)

type RotatingLogger struct {
	mu       sync.Mutex
	file     *os.File
	size     int64
	basePath string
	sequence int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file == nil {
		return fmt.Errorf("no open file")
	}

	rl.file.Close()

	for i := maxBackups - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	os.Rename(rl.basePath, rl.backupPath(0))
	return rl.openCurrent()
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.basePath
	}
	ext := filepath.Ext(rl.basePath)
	base := rl.basePath[:len(rl.basePath)-len(ext)]
	return fmt.Sprintf("%s.%d%s", base, index, ext)
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		rl.mu.Unlock()
		rl.rotate()
		rl.mu.Lock()
	}

	n, err = rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) WriteString(s string) (n int, err error) {
	return rl.Write([]byte(s))
}

func (rl *RotatingLogger) Log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)
	rl.WriteString(logLine)
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.file.Close()

	for i := 0; i < 1000; i++ {
		logger.Log("INFO", fmt.Sprintf("Log entry number %d", i))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}