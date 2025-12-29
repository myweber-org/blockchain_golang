package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logFileName    = "app.log"
)

type LogRotator struct {
	currentFile *os.File
	filePath    string
	baseName    string
	dir         string
}

func NewLogRotator() (*LogRotator, error) {
	dir := "./logs"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(dir, logFileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &LogRotator{
		currentFile: file,
		filePath:    filePath,
		baseName:    strings.TrimSuffix(logFileName, filepath.Ext(logFileName)),
		dir:         dir,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	stat, err := lr.currentFile.Stat()
	if err != nil {
		return 0, err
	}

	if stat.Size()+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	return lr.currentFile.Write(p)
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(lr.dir, fmt.Sprintf("%s_%s.log", lr.baseName, timestamp))

	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	lr.currentFile = file

	go lr.cleanupOldFiles()

	return nil
}

func (lr *LogRotator) cleanupOldFiles() {
	files, err := filepath.Glob(filepath.Join(lr.dir, lr.baseName+"_*.log"))
	if err != nil {
		return
	}

	if len(files) <= maxBackupFiles {
		return
	}

	sort.Strings(files)
	filesToDelete := files[:len(files)-maxBackupFiles]

	for _, file := range filesToDelete {
		os.Remove(file)
	}
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Iteration %d: Log message here\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}