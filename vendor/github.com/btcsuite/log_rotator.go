package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxBackupCount = 5
	logFileName   = "app.log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator(baseDir string) (*LogRotator, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(baseDir, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		basePath:    baseDir,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
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
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	backupPath := filepath.Join(lr.basePath, backupName)

	if err := os.Rename(filepath.Join(lr.basePath, logFileName), backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(lr.basePath, logFileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0

	go lr.cleanupOldBackups()

	return nil
}

func (lr *LogRotator) cleanupOldBackups() {
	files, err := filepath.Glob(filepath.Join(lr.basePath, logFileName+".*"))
	if err != nil {
		return
	}

	if len(files) <= maxBackupCount {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return extractTimestamp(files[i]) > extractTimestamp(files[j])
	})

	for i := maxBackupCount; i < len(files); i++ {
		os.Remove(files[i])
	}
}

func extractTimestamp(path string) time.Time {
	base := filepath.Base(path)
	parts := strings.Split(base, ".")
	if len(parts) < 2 {
		return time.Time{}
	}

	t, err := time.Parse("20060102_150405", parts[1])
	if err != nil {
		return time.Time{}
	}
	return t
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("./logs")
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	writer := io.MultiWriter(os.Stdout, rotator)

	for i := 0; i < 1000; i++ {
		fmt.Fprintf(writer, "Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}
}