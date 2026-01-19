
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	dateStr := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != dateStr {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLog(); err != nil {
				return err
			}
			if err := rl.cleanupOldBackups(); err != nil {
				return err
			}
		}

		filename := fmt.Sprintf("%s.%s.log", rl.basePath, dateStr)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
		rl.currentDay = dateStr
	}
	return nil
}

func (rl *RotatingLogger) compressOldLog() error {
	oldPath := fmt.Sprintf("%s.%s.log", rl.basePath, rl.currentDay)
	compressedPath := oldPath + ".gz"

	src, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	pattern := rl.basePath + ".*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, file := range toDelete {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))))
		time.Sleep(100 * time.Millisecond)
	}
}