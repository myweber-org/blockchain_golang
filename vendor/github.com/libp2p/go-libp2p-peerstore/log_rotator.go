package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	basePath    string
	fileIndex   int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:  strings.TrimSuffix(path, filepath.Ext(path)),
		fileIndex: 0,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
		if err := compressFile(rl.currentFile.Name()); err != nil {
			fmt.Printf("Compression failed: %v\n", err)
		}
	}

	rl.fileIndex++
	if rl.fileIndex > maxBackups {
		rl.fileIndex = 1
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := fmt.Sprintf("%s.%d.log", rl.basePath, rl.fileIndex)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func compressFile(source string) error {
	dest := source + ".gz"
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Simple copy for demonstration (replace with actual compression)
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		os.Remove(dest)
		return err
	}

	os.Remove(source)
	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}