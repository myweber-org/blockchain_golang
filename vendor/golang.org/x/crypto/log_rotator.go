package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logFileName = "app.log"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	fullPath := filepath.Join(path, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		currentFile: file,
		currentSize: info.Size(),
		basePath:    path,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	backupPath := filepath.Join(rl.basePath, backupName)

	if err := os.Rename(filepath.Join(rl.basePath, logFileName), backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(rl.basePath, logFileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0

	go rl.cleanupOldBackups()

	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := filepath.Join(rl.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackups {
		return
	}

	toDelete := matches[:len(matches)-maxBackups]
	for _, file := range toDelete {
		os.Remove(file)
	}
}

func (rl *RotatingLogger) Close() error {
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger(".")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}