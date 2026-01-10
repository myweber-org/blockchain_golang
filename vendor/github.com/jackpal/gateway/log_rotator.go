
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
	basePath     string
	currentFile  *os.File
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     maxSize,
		backupCount: backupCount,
		compressOld: compressOld,
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

	n, err := rl.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentFile == nil || rl.currentSize >= rl.maxSize {
		return rl.performRotation()
	}
	return nil
}

func (rl *RotatingLogger) performRotation() error {
	if rl.currentFile != nil {
		if err := rl.currentFile.Close(); err != nil {
			return err
		}

		if err := rl.rotateExistingFiles(); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) rotateExistingFiles() error {
	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)

	for i := rl.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string

		if i == 0 {
			oldPath = rl.basePath
		} else {
			oldPath = filepath.Join(dir, fmt.Sprintf("%s.%d", baseName, i))
			if rl.compressOld {
				oldPath += ".gz"
			}
		}

		newPath = filepath.Join(dir, fmt.Sprintf("%s.%d", baseName, i+1))
		if rl.compressOld {
			newPath += ".gz"
		}

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if rl.compressOld {
		return rl.compressFile(rl.basePath + ".1")
	}
	return nil
}

func (rl *RotatingLogger) compressFile(path string) error {
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
	logger, err := NewRotatingLogger("app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}