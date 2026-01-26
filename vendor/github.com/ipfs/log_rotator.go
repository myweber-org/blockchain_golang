package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	filePath      string
	maxSize       int64
	currentSize   int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}
	
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	if rl.currentSize+int64(len(p)) > rl.maxSize {
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
	rl.currentFile.Close()
	
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.bak", rl.filePath, timestamp)
	
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}
	
	rl.rotationCount++
	if err := rl.openCurrentFile(); err != nil {
		return err
	}
	
	go rl.compressBackup(backupPath)
	return nil
}

func (rl *RotatingLogger) compressBackup(backupPath string) {
	compressedPath := backupPath + ".gz"
	
	src, err := os.Open(backupPath)
	if err != nil {
		return
	}
	defer src.Close()
	
	dst, err := os.Create(compressedPath)
	if err != nil {
		return
	}
	defer dst.Close()
	
	gzWriter := gzip.NewWriter(dst)
	defer gzWriter.Close()
	
	if _, err := io.Copy(gzWriter, src); err == nil {
		os.Remove(backupPath)
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()
	
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Printf("Log rotation completed %d times\n", logger.rotationCount)
}