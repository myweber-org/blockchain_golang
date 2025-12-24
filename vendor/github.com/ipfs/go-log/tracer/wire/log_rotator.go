
package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const maxLogSize = 10 * 1024 * 1024 // 10MB

type LogRotator struct {
	filePath string
	file     *os.File
	currentSize int64
}

func NewLogRotator(path string) (*LogRotator, error) {
	rotator := &LogRotator{filePath: path}
	err := rotator.openCurrentFile()
	return rotator, err
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxLogSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}
	
	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if lr.file != nil {
		lr.file.Close()
	}
	
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	backupPath := lr.filePath + "." + timestamp
	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}
	
	return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
	dir := filepath.Dir(lr.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	
	lr.file = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}