package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	basePath   string
	maxSize    int64
	currentSize int64
}

func NewRotatingLogger(basePath string, maxSize int64) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()
	for i := 9; i >= 1; i-- {
		oldName := fmt.Sprintf("%s.%d", rl.basePath, i)
		newName := fmt.Sprintf("%s.%d", rl.basePath, i+1)
		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}
	backupName := fmt.Sprintf("%s.1", rl.basePath)
	os.Rename(rl.basePath, backupName)
	return rl.openCurrent()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 1024*1024)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d: Some sample log data here.\n", i)
		logger.Write([]byte(msg))
	}
	fmt.Println("Log rotation test completed. Check app.log and rotated files.")
}