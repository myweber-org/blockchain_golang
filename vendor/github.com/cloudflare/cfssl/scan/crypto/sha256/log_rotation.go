package main

import (
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
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	baseName := filepath.Join(logDir, name)
	rl := &RotatingLogger{baseName: baseName}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := rl.baseName + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	timestamp := time.Now().Format("20060102_150405")
	rotatedName := rl.baseName + "_" + timestamp + ".log"
	if err := os.Rename(rl.baseName+".log", rotatedName); err != nil {
		return err
	}

	if err := rl.cleanupOldFiles(); err != nil {
		log.Printf("Cleanup error: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := rl.baseName + "_*.log"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackups {
		return nil
	}

	toDelete := len(matches) - maxBackups
	for i := 0; i < toDelete; i++ {
		if err := os.Remove(matches[i]); err != nil {
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

func extractTimestamp(filename string) (time.Time, error) {
	base := filepath.Base(filename)
	parts := strings.Split(base, "_")
	if len(parts) < 3 {
		return time.Time{}, os.ErrInvalid
	}

	timestampStr := parts[1] + "_" + strings.TrimSuffix(parts[2], ".log")
	return time.Parse("20060102_150405", timestampStr)
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := "Log entry " + strconv.Itoa(i) + "\n"
		if _, err := logger.Write([]byte(message)); err != nil {
			log.Printf("Write error: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}