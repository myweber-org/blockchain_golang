package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	
	dir := filepath.Dir(basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logger := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.currentFile.Close()
	
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.%d", rl.filePath, timestamp, rl.rotationCount)
	
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	rl.rotationCount++
	return rl.openCurrentFile()
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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
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
	mu            sync.Mutex
	basePath      string
	currentFile   *os.File
	maxSize       int64
	currentSize   int64
	backupCount   int
	compressOld   bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:    basePath,
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

	file, err := os.OpenFile(rl.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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

	// Rename existing backups
	for i := rl.backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	// Rename current file to backup.0
	backupPath := rl.getBackupPath(0)
	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return err
	}

	// Compress if enabled
	if rl.compressOld {
		go rl.compressFile(backupPath)
	}

	// Open new current file
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d", rl.basePath, index+1)
}

func (rl *RotatingLogger) compressFile(path string) {
	// Simulated compression - in real implementation use gzip or similar
	// This is just a placeholder to show the concept
	if strings.HasSuffix(path, ".log") {
		compressedPath := path + ".gz"
		// Actual compression logic would go here
		fmt.Printf("Compressing %s to %s\n", path, compressedPath)
		os.Remove(path)
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Use the rotating logger
	customLog := log.New(logger, "APP: ", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d at %v", i, time.Now())
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
}