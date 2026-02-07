package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type RotatingLogger struct {
	mu          sync.Mutex
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
}

func NewRotatingLogger(path string, maxSizeMB int) (*RotatingLogger, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		filePath: absPath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	info, err := os.Stat(rl.filePath)
	if err == nil {
		rl.currentSize = info.Size()
	} else if os.IsNotExist(err) {
		rl.currentSize = 0
	} else {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.file = file
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	for i := 9; i > 0; i-- {
		oldName := fmt.Sprintf("%s.%d", rl.filePath, i)
		newName := fmt.Sprintf("%s.%d", rl.filePath, i+1)
		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	backupName := fmt.Sprintf("%s.1", rl.filePath)
	os.Rename(rl.filePath, backupName)

	return rl.openFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
		rl.currentSize = 0
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
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		logger.Write([]byte(msg))
	}

	fmt.Println("Log rotation test completed")
}