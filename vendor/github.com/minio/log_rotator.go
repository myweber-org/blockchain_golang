package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
		sequence: 0,
	}

	if err := rl.openNewFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openNewFile() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
		rl.compressOldFile()
	}

	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0
	rl.sequence++
	return nil
}

func (rl *RotatingLogger) compressOldFile() {
	if rl.sequence <= 1 {
		return
	}

	oldFilename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence-2))
	newFilename := oldFilename + ".gz"

	oldFile, err := os.Open(oldFilename)
	if err != nil {
		return
	}
	defer oldFile.Close()

	newFile, err := os.Create(newFilename)
	if err != nil {
		return
	}
	defer newFile.Close()

	gzWriter := gzip.NewWriter(newFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		return
	}

	os.Remove(oldFilename)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.openNewFile(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
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
	fileIndex    int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	if err := rl.openOrCreateLog(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openOrCreateLog() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	ext := filepath.Ext(rl.basePath)
	base := rl.basePath[:len(rl.basePath)-len(ext)]
	path := fmt.Sprintf("%s_%d%s", base, rl.fileIndex, ext)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		rl.fileIndex++
		if err := rl.openOrCreateLog(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
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
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}