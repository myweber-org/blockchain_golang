
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
	maxFileSize  = 10 * 1024 * 1024
	maxBackups   = 5
	logDirectory = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
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
	filename := filepath.Join(logDirectory, rl.baseName+".log")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Failed to rotate log file: %v", err)
		}
	}

	n, err := rl.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		if err := rl.currentFile.Close(); err != nil {
			return fmt.Errorf("failed to close current log file: %w", err)
		}
	}

	timestamp := time.Now().Format("20060102-150405")
	oldPath := filepath.Join(logDirectory, rl.baseName+".log")
	newPath := filepath.Join(logDirectory, fmt.Sprintf("%s-%s.log", rl.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if err := rl.compressFile(newPath); err != nil {
		log.Printf("Failed to compress log file %s: %v", newPath, err)
	}

	if err := rl.cleanupOldBackups(); err != nil {
		log.Printf("Failed to cleanup old backups: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(source + ".gz")
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	if err := os.Remove(source); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	pattern := filepath.Join(logDirectory, rl.baseName+"-*.log.gz")
	backups, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(backups) <= maxBackups {
		return nil
	}

	for i := 0; i < len(backups)-maxBackups; i++ {
		if err := os.Remove(backups[i]); err != nil {
			return err
		}
	}

	return nil
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
	logger, err := NewRotatingLogger("application")
	if err != nil {
		log.Fatalf("Failed to create rotating logger: %v", err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: This is a test log message for rotation testing", i)
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	logDirectory  string
	baseFilename  string
	retentionDays int
}

func NewRotatingLogger(dir, filename string, maxSize int64, retention int) (*RotatingLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fullPath := filepath.Join(dir, filename)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &RotatingLogger{
		currentFile:   file,
		currentSize:   info.Size(),
		maxSize:       maxSize,
		logDirectory:  dir,
		baseFilename:  filename,
		retentionDays: retention,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("rotation failed: %w", err)
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
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(rl.logDirectory, rl.baseFilename)
	newPath := filepath.Join(rl.logDirectory, fmt.Sprintf("%s.%s.log", rl.baseFilename, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if err := rl.compressFile(newPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to compress %s: %v\n", newPath, err)
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	rl.currentFile = file
	rl.currentSize = 0

	go rl.cleanupOldLogs()
	return nil
}

func (rl *RotatingLogger) compressFile(source string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(source + ".gz")
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	if err := os.Remove(source); err != nil {
		return fmt.Errorf("failed to remove uncompressed file: %w", err)
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
	cutoffTime := time.Now().AddDate(0, 0, -rl.retentionDays)
	pattern := filepath.Join(rl.logDirectory, rl.baseFilename+".*.log.gz")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	for _, file := range matches {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			os.Remove(file)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp", "application.log", 10*1024*1024, 30)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))))
		time.Sleep(10 * time.Millisecond)
	}
}