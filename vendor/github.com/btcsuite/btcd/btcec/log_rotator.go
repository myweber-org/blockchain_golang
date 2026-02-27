
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
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	fileCounter  int
	maxFiles     int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int, compressOld bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     maxSize,
		maxFiles:    maxFiles,
		compressOld: compressOld,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	ext := filepath.Ext(rl.basePath)
	base := strings.TrimSuffix(rl.basePath, ext)

	files, _ := filepath.Glob(base + "_*" + ext)
	rl.fileCounter = len(files)

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
		rl.currentFile = nil
	}

	ext := filepath.Ext(rl.basePath)
	base := strings.TrimSuffix(rl.basePath, ext)

	timestamp := time.Now().Format("20060102_150405")
	newPath := fmt.Sprintf("%s_%s%s", base, timestamp, ext)

	if err := os.Rename(rl.basePath, newPath); err != nil {
		return err
	}

	rl.fileCounter++

	if rl.compressOld {
		go rl.compressFile(newPath)
	}

	if rl.fileCounter > rl.maxFiles {
		go rl.cleanupOldFiles(base, ext)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(path string) {
	// Compression implementation would go here
	// For now just log that compression would happen
	log.Printf("Would compress file: %s", path)
}

func (rl *RotatingLogger) cleanupOldFiles(base, ext string) {
	pattern := fmt.Sprintf("%s_*%s", base, ext)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(files) > rl.maxFiles {
		filesToDelete := files[:len(files)-rl.maxFiles]
		for _, file := range filesToDelete {
			os.Remove(file)
			log.Printf("Removed old log file: %s", file)
		}
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
	logger, err := NewRotatingLogger("app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d: This is a test log message that will trigger rotation", i)
		time.Sleep(100 * time.Millisecond)
	}
}