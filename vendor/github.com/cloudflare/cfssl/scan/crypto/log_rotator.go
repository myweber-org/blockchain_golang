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

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	sequence    int
}

func NewLogRotator(baseName string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	lr := &LogRotator{
		baseName: baseName,
		sequence: 0,
	}

	if err := lr.openNewFile(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openNewFile() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", lr.baseName, lr.sequence))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0
	lr.sequence++

	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	lr.currentSize += int64(n)
	return n, err
}

func (lr *LogRotator) rotate() error {
	oldFile := lr.currentFile
	oldPath := oldFile.Name()

	if err := lr.openNewFile(); err != nil {
		return err
	}

	go func() {
		compressedPath := oldPath + ".gz"
		if err := compressFile(oldPath, compressedPath); err == nil {
			os.Remove(oldPath)
		}
	}()

	return nil
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

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (lr *LogRotator) Close() error {
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
		logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		rotator.Write([]byte(logEntry))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}