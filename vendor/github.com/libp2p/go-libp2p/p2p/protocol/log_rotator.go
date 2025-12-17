
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
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewLogRotator(baseName string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	lr := &LogRotator{
		baseName: baseName,
	}

	if err := lr.openCurrentFile(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", lr.baseName, timestamp))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	files, err := filepath.Glob(filepath.Join(logDir, lr.baseName+"_*.log"))
	if err != nil {
		return err
	}

	if len(files) >= maxBackups {
		oldest := files[0]
		for _, f := range files[1:] {
			info1, _ := os.Stat(oldest)
			info2, _ := os.Stat(f)
			if info2.ModTime().Before(info1.ModTime()) {
				oldest = f
			}
		}
		if err := compressAndRemove(oldest); err != nil {
			return err
		}
	}

	return lr.openCurrentFile()
}

func compressAndRemove(filename string) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(filename)
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app")
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		rotator.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}