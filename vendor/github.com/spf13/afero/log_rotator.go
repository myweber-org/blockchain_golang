package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	baseName    string
}

func NewLogRotator(filename string) (*LogRotator, error) {
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

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		baseName:    filename,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(logDir, lr.baseName)
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.%s", lr.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0

	go lr.cleanupOldFiles()
	return nil
}

func (lr *LogRotator) cleanupOldFiles() {
	pattern := filepath.Join(logDir, fmt.Sprintf("%s.*", lr.baseName))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackups {
		return
	}

	filesToDelete := matches[:len(matches)-maxBackups]
	for _, file := range filesToDelete {
		os.Remove(file)
	}
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("app.log")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}