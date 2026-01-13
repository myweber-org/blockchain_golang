
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
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
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
	if err != nil {
		return n, err
	}
	rl.size += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	today := time.Now().Format("2006-01-02")
	if rl.currentDay != today {
		if err := rl.rotate(); err != nil {
			return err
		}
		rl.currentDay = today
	}

	if rl.size >= maxFileSize {
		return rl.rotate()
	}

	if rl.file == nil {
		return rl.openFile()
	}
	return nil
}

func (rl *RotatingLogger) openFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := rl.basePath + ".log"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	if info, err := file.Stat(); err == nil {
		rl.size = info.Size()
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
		rl.file = nil
	}

	oldPath := rl.basePath + ".log"
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return rl.openFile()
	}

	timestamp := time.Now().Format("20060102-150405")
	newPath := fmt.Sprintf("%s.%s.log", rl.basePath, timestamp)
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	go compressLogFile(newPath)
	go cleanupOldBackups(rl.basePath)

	return rl.openFile()
}

func compressLogFile(path string) {
	compressedPath := path + ".gz"
	inFile, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open log file for compression: %v", err)
		return
	}
	defer inFile.Close()

	outFile, err := os.Create(compressedPath)
	if err != nil {
		log.Printf("Failed to create compressed file: %v", err)
		return
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, inFile); err != nil {
		log.Printf("Failed to compress log file: %v", err)
		return
	}

	os.Remove(path)
}

func cleanupOldBackups(basePath string) {
	pattern := basePath + ".*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		os.Remove(matches[i])
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
	logger, err := NewRotatingLogger("/var/log/myapp/app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Application is running normally", i)
		time.Sleep(100 * time.Millisecond)
	}
}