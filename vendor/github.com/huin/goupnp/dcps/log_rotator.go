
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
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu         sync.Mutex
	current    *os.File
	size       int64
	baseName   string
	fileNumber int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}
	if err := rl.openNewFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
	rl.fileNumber++
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.fileNumber))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if rl.current != nil {
		rl.current.Close()
		go rl.compressPreviousFile(rl.fileNumber - 1)
	}
	rl.current = file
	rl.size = 0
	return nil
}

func (rl *RotatingLogger) compressPreviousFile(number int) {
	oldFile := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, number))
	compressedFile := oldFile + ".gz"

	src, err := os.Open(oldFile)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(compressedFile)
	if err != nil {
		return
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return
	}
	os.Remove(oldFile)
	rl.cleanupOldBackups()
}

func (rl *RotatingLogger) cleanupOldBackups() {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log.gz"))
	if err != nil {
		return
	}
	if len(files) > maxBackups {
		for i := 0; i < len(files)-maxBackups; i++ {
			os.Remove(files[i])
		}
	}
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.openNewFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.current.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.current != nil {
		return rl.current.Close()
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
		msg := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
}