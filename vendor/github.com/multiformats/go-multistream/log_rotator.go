package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
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

	oldPath := filepath.Join(logDir, rl.baseName+".log")
	newPath := filepath.Join(logDir, rl.baseName+"."+time.Now().Format("20060102150405")+".log")

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.compressOldLogs(); err != nil {
		log.Printf("Failed to compress logs: %v", err)
	}

	if err := rl.cleanupOldLogs(); err != nil {
		log.Printf("Failed to cleanup old logs: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressOldLogs() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+".*.log"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".gz") {
			continue
		}

		if err := compressFile(file); err != nil {
			return err
		}
	}

	return nil
}

func compressFile(src string) error {
	dest := src + ".gz"
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Simple copy for demonstration (in real implementation use gzip.Writer)
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		os.Remove(dest)
		return err
	}

	if err := os.Remove(src); err != nil {
		os.Remove(dest)
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldLogs() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+".*.gz"))
	if err != nil {
		return err
	}

	if len(files) <= maxBackups {
		return nil
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(files)-maxBackups; i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}

	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}