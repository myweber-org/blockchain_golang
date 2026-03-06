
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	basePath     string
	currentFile  *os.File
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		basePath:    absPath,
		maxSize:     int64(maxSizeMB) * 1024 * 1024,
		backupCount: backupCount,
		compressOld: compressOld,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return err
	}

	if rl.compressOld {
		go rl.compressFile(backupPath)
	}

	rl.cleanOldBackups()

	if err := rl.openCurrentFile(); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) compressFile(path string) {
	// Compression logic would go here
	// For simplicity, we just log the compression event
	log.Printf("Compressing file: %s", path)
	// In real implementation, use compress/gzip or similar
}

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := fmt.Sprintf("%s.*", filepath.Base(rl.basePath))
	dir := filepath.Dir(rl.basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var backups []string
	for _, entry := range entries {
		if matched, _ := filepath.Match(pattern, entry.Name()); matched {
			backups = append(backups, filepath.Join(dir, entry.Name()))
		}
	}

	if len(backups) <= rl.backupCount {
		return
	}

	// Sort by modification time (oldest first)
	// For simplicity, we just remove the first N files
	// In production, implement proper sorting
	for i := 0; i < len(backups)-rl.backupCount; i++ {
		os.Remove(backups[i])
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Example usage
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
}