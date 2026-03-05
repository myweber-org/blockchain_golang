
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	baseName   string
	currentDay string
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, filename)
	rl := &RotatingLogger{
		baseName: basePath,
	}

	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDate {
		if err := rl.rotate(); err != nil {
			return err
		}
		rl.currentDay = currentDate
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
		if err := rl.compressOldLog(); err != nil {
			log.Printf("Failed to compress log: %v", err)
		}
		rl.cleanupOldBackups()
	}

	newPath := rl.baseName + "." + time.Now().Format("2006-01-02_150405")
	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.size = 0

	// Create symlink to current log
	symlinkPath := rl.baseName + ".current"
	os.Remove(symlinkPath)
	os.Symlink(filepath.Base(newPath), symlinkPath)

	return nil
}

func (rl *RotatingLogger) compressOldLog() error {
	logFiles, err := filepath.Glob(rl.baseName + ".*")
	if err != nil {
		return err
	}

	for _, file := range logFiles {
		if filepath.Ext(file) == ".gz" {
			continue
		}

		if err := compressFile(file); err != nil {
			return err
		}
		os.Remove(file)
	}
	return nil
}

func compressFile(src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(src + ".gz")
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
}

func (rl *RotatingLogger) cleanupOldBackups() {
	files, err := filepath.Glob(rl.baseName + ".*.gz")
	if err != nil {
		return
	}

	if len(files) > backupCount {
		filesToDelete := files[:len(files)-backupCount]
		for _, file := range filesToDelete {
			os.Remove(file)
		}
	}
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
	logger, err := NewRotatingLogger("application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Simulate log writing
	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Processing request from user %d\n",
			time.Now().Format(time.RFC3339), i, i%100)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed. Check ./logs directory")
}