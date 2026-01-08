
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
	if maxSizeMB <= 0 {
		return nil, fmt.Errorf("maxSizeMB must be positive")
	}
	if backupCount < 0 {
		return nil, fmt.Errorf("backupCount cannot be negative")
	}

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     int64(maxSizeMB) * 1024 * 1024,
		backupCount: backupCount,
		compressOld: compressOld,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
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
	}

	ext := filepath.Ext(rl.basePath)
	base := strings.TrimSuffix(rl.basePath, ext)
	timestamp := time.Now().Format("20060102_150405")

	for i := rl.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		if i == 0 {
			oldPath = rl.basePath
		} else {
			oldPath = fmt.Sprintf("%s.%d%s", base, i, ext)
		}

		if rl.compressOld && i == rl.backupCount-1 && rl.backupCount > 0 {
			newPath = fmt.Sprintf("%s.%s%s.gz", base, timestamp, ext)
		} else if i+1 < rl.backupCount {
			newPath = fmt.Sprintf("%s.%d%s", base, i+1, ext)
		} else {
			newPath = fmt.Sprintf("%s.%s%s", base, timestamp, ext)
		}

		if _, err := os.Stat(oldPath); err == nil {
			if rl.compressOld && i == rl.backupCount-1 && rl.backupCount > 0 {
				if err := compressFile(oldPath, newPath); err != nil {
					log.Printf("Failed to compress %s: %v", oldPath, err)
					os.Rename(oldPath, newPath)
				}
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	return rl.openCurrentFile()
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

	compressor := NewGzipWriter(dstFile)
	defer compressor.Close()

	_, err = io.Copy(compressor, srcFile)
	if err != nil {
		os.Remove(dst)
		return err
	}

	os.Remove(src)
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

type GzipWriter struct {
	io.WriteCloser
}

func NewGzipWriter(w io.Writer) *GzipWriter {
	return &GzipWriter{}
}

func (gw *GzipWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (gw *GzipWriter) Close() error {
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation demonstration completed")
}