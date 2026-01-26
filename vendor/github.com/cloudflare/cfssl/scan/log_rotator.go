
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
	backupCount = 5
)

type RotatingFile struct {
	mu         sync.Mutex
	filename   string
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingFile(filename string) (*RotatingFile, error) {
	rf := &RotatingFile{
		filename: filename,
		basePath: filepath.Dir(filename),
	}

	if err := rf.openCurrent(); err != nil {
		return nil, err
	}

	return rf, nil
}

func (rf *RotatingFile) openCurrent() error {
	info, err := os.Stat(rf.filename)
	if os.IsNotExist(err) {
		file, err := os.Create(rf.filename)
		if err != nil {
			return err
		}
		rf.file = file
		rf.size = 0
		return nil
	}
	if err != nil {
		return err
	}

	file, err := os.OpenFile(rf.filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rf.file = file
	rf.size = info.Size()
	return nil
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.size+int64(len(p)) > maxFileSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.file.Write(p)
	if err == nil {
		rf.size += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) rotate() error {
	if err := rf.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s.gz", rf.filename, timestamp)

	if err := compressFile(rf.filename, backupName); err != nil {
		return err
	}

	if err := os.Remove(rf.filename); err != nil {
		return err
	}

	rf.cleanOldBackups()

	return rf.openCurrent()
}

func compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (rf *RotatingFile) cleanOldBackups() {
	pattern := rf.filename + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		os.Remove(matches[i])
	}
}

func (rf *RotatingFile) Close() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.file.Close()
}

func main() {
	logFile, err := NewRotatingFile("application.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logFile.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}