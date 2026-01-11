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
}package logutil

import (
	"compress/gzip"
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
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentSize+int64(len(p)) > rl.maxSize && rl.currentSize > 0 {
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
		
		oldPath := rl.currentFile.Name()
		timestamp := time.Now().Format("20060102_150405")
		archivePath := fmt.Sprintf("%s.%s.gz", oldPath, timestamp)
		
		if err := compressFile(oldPath, archivePath); err != nil {
			return err
		}
		
		os.Remove(oldPath)
	}
	
	rl.rotationCount++
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	f, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	
	rl.currentFile = f
	rl.currentSize = info.Size()
	return nil
}

func compressFile(source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	
	gz := gzip.NewWriter(out)
	defer gz.Close()
	
	_, err = io.Copy(gz, in)
	return err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}