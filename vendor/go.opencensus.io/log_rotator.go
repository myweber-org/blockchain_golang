
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

const (
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	backupCount    = 5
	compressBackup = true
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
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	if rl.file == nil || rl.currentDay != currentDate || rl.size >= maxFileSize {
		if rl.file != nil {
			rl.file.Close()
		}

		rl.currentDay = currentDate
		newPath := rl.generateFilename()

		if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		info, err := file.Stat()
		if err != nil {
			file.Close()
			return err
		}

		rl.file = file
		rl.size = info.Size()

		go rl.cleanupOldFiles()
	}
	return nil
}

func (rl *RotatingLogger) generateFilename() string {
	base := strings.TrimSuffix(rl.basePath, filepath.Ext(rl.basePath))
	ext := filepath.Ext(rl.basePath)
	if ext == "" {
		ext = ".log"
	}
	return fmt.Sprintf("%s-%s%s", base, rl.currentDay, ext)
}

func (rl *RotatingLogger) cleanupOldFiles() {
	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to read directory: %v", err)
		return
	}

	var logFiles []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), baseName) && strings.HasSuffix(file.Name(), ".log") {
			logFiles = append(logFiles, filepath.Join(dir, file.Name()))
		}
	}

	if len(logFiles) > backupCount {
		filesToRemove := logFiles[:len(logFiles)-backupCount]
		for _, file := range filesToRemove {
			if compressBackup {
				rl.compressFile(file)
			} else {
				os.Remove(file)
			}
		}
	}
}

func (rl *RotatingLogger) compressFile(path string) {
	src, err := os.Open(path)
	if err != nil {
		return
	}
	defer src.Close()

	dstPath := path + ".gz"
	dst, err := os.Create(dstPath)
	if err != nil {
		return
	}
	defer dst.Close()

	// In production, use compress/gzip here
	// This is simplified for demonstration
	_, err = io.Copy(dst, src)
	if err == nil {
		os.Remove(path)
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
	logger, err := NewRotatingLogger("./logs/application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "APP: ", log.LstdFlags)
	for i := 0; i < 100; i++ {
		customLog.Printf("Log entry number %d", i)
		time.Sleep(100 * time.Millisecond)
	}
}