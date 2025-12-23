package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRotator struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	backupCount int
	currentSize int64
}

func NewLogRotator(filePath string, maxSizeMB int, backupCount int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rotator := &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	dir := filepath.Dir(lr.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %w", err)
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open log file failed: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("stat log file failed: %w", err)
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("rotate failed: %w", err)
		}
	}

	n, err := lr.currentFile.Write(p)
	if err != nil {
		return n, fmt.Errorf("write to log file failed: %w", err)
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return fmt.Errorf("close current file failed: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return fmt.Errorf("rename log file failed: %w", err)
	}

	if err := lr.compressBackup(backupPath); err != nil {
		fmt.Fprintf(os.Stderr, "compress backup failed: %v\n", err)
	}

	if err := lr.cleanupOldBackups(); err != nil {
		fmt.Fprintf(os.Stderr, "cleanup old backups failed: %v\n", err)
	}

	if err := lr.openCurrentFile(); err != nil {
		return fmt.Errorf("reopen log file failed: %w", err)
	}

	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) compressBackup(srcPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open backup file failed: %w", err)
	}
	defer srcFile.Close()

	dstPath := srcPath + ".gz"
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create compressed file failed: %w", err)
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return fmt.Errorf("compress data failed: %w", err)
	}

	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("remove uncompressed backup failed: %w", err)
	}

	return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
	pattern := lr.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob backup files failed: %w", err)
	}

	if len(matches) <= lr.backupCount {
		return nil
	}

	backups := make([]string, len(matches))
	copy(backups, matches)

	for i := 0; i < len(backups)-lr.backupCount; i++ {
		if err := os.Remove(backups[i]); err != nil {
			return fmt.Errorf("remove old backup %s failed: %w", backups[i], err)
		}
	}

	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create log rotator failed: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Fprintf(os.Stderr, "write log failed: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}