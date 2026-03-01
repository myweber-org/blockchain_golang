
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
	maxLogSize    = 10 * 1024 * 1024 // 10MB
	backupCount   = 5
	logDir        = "./logs"
	currentLog    = "app.log"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	basePath  string
}

func NewRotatingLogger() (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	path := filepath.Join(logDir, currentLog)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	info, _ := file.Stat()
	return &RotatingLogger{
		file:     file,
		size:     info.Size(),
		basePath: path,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxLogSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", timestamp))

	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.size = 0

	go compressOldLog(backupPath)
	go cleanupOldBackups()

	return nil
}

func compressOldLog(path string) {
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

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return
	}

	os.Remove(path)
}

func cleanupOldBackups() {
	pattern := filepath.Join(logDir, "app_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		os.Remove(matches[i])
	}
}

func main() {
	logger, err := NewRotatingLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.file.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Application is running normally", i)
		time.Sleep(10 * time.Millisecond)
	}
}