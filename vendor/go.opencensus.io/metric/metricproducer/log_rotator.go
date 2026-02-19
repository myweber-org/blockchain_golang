package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type LogRotator struct {
	FilePath    string
	MaxSize     int64
	MaxBackups  int
	currentSize int64
}

func NewLogRotator(filePath string, maxSize int64, maxBackups int) *LogRotator {
	return &LogRotator{
		FilePath:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
	}
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
	if lr.currentSize+int64(len(p)) > lr.MaxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	file, err := os.OpenFile(lr.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err = file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	for i := lr.MaxBackups - 1; i > 0; i-- {
		oldPath := lr.FilePath + "." + strconv.Itoa(i)
		newPath := lr.FilePath + "." + strconv.Itoa(i+1)
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	backupPath := lr.FilePath + ".1"
	if _, err := os.Stat(lr.FilePath); err == nil {
		if err := os.Rename(lr.FilePath, backupPath); err != nil {
			return err
		}
	}

	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) CleanOldLogs() error {
	for i := lr.MaxBackups + 1; i <= 10; i++ {
		path := lr.FilePath + "." + strconv.Itoa(i)
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}
	return nil
}

func main() {
	rotator := NewLogRotator("app.log", 1024*1024, 5)
	
	logData := []byte(time.Now().Format("2006-01-02 15:04:05") + " - Application started\n")
	rotator.Write(logData)
	
	rotator.CleanOldLogs()
}