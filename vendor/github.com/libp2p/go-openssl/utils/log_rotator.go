
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
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) (*LogRotator, error) {
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
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	lr.currentFile = file
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	stat, err := lr.currentFile.Stat()
	if err != nil {
		return 0, err
	}

	if stat.Size()+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	return lr.currentFile.Write(p)
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	if err := lr.compressFile(backupPath); err != nil {
		return err
	}

	if err := lr.cleanupOldBackups(); err != nil {
		return err
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(sourcePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	compressedPath := sourcePath + ".gz"
	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	gzipWriter := gzip.NewWriter(compressedFile)
	defer gzipWriter.Close()

	if _, err := io.Copy(gzipWriter, sourceFile); err != nil {
		return err
	}

	if err := os.Remove(sourcePath); err != nil {
		return err
	}

	return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
	pattern := lr.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= lr.backupCount {
		return nil
	}

	backupsToRemove := matches[:len(matches)-lr.backupCount]
	for _, backup := range backupsToRemove {
		if err := os.Remove(backup); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 10*1024*1024, 5)
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application event occurred at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}