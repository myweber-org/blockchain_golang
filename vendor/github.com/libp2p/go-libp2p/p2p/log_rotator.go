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
	mu         sync.Mutex
	file       *os.File
	basePath   string
	maxSize    int64
	fileSize   int64
	backupCount int
}

func NewRotatingLogger(basePath string, maxSize int64, backupCount int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:   basePath,
		maxSize:    maxSize,
		backupCount: backupCount,
	}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.fileSize += int64(n)
	if rl.fileSize >= rl.maxSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)
	if err := compressFile(rl.basePath, backupPath); err != nil {
		return err
	}
	if err := os.Remove(rl.basePath); err != nil {
		return err
	}
	if err := rl.openFile(); err != nil {
		return err
	}
	return rl.cleanupOldBackups()
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	gzWriter := newGzipWriter(dstFile)
	defer gzWriter.Close()
	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) <= rl.backupCount {
		return nil
	}
	toDelete := matches[:len(matches)-rl.backupCount]
	for _, path := range toDelete {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10*1024*1024, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
	log.SetOutput(logger)
	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("x", 1024))
		time.Sleep(10 * time.Millisecond)
	}
}