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

type LogRotator struct {
	filePath    string
	maxSize     int64
	maxFiles    int
	currentSize int64
	file        *os.File
}

func NewLogRotator(filePath string, maxSizeMB int, maxFiles int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
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

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.file.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return err
	}

	if err := lr.renameOldFiles(); err != nil {
		return err
	}

	if err := os.Rename(lr.filePath, lr.filePath+".1"); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := lr.openCurrentFile(); err != nil {
		return err
	}

	return lr.cleanupOldFiles()
}

func (lr *LogRotator) renameOldFiles() error {
	dir := filepath.Dir(lr.filePath)
	base := filepath.Base(lr.filePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var oldFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, base+".") && entry.Type().IsRegular() {
			oldFiles = append(oldFiles, name)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(oldFiles)))

	for _, oldFile := range oldFiles {
		parts := strings.Split(oldFile, ".")
		if len(parts) < 2 {
			continue
		}

		index, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			continue
		}

		newName := base + "." + strconv.Itoa(index+1)
		oldPath := filepath.Join(dir, oldFile)
		newPath := filepath.Join(dir, newName)

		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) cleanupOldFiles() error {
	dir := filepath.Dir(lr.filePath)
	base := filepath.Base(lr.filePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var filesToDelete []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, base+".") && entry.Type().IsRegular() {
			parts := strings.Split(name, ".")
			if len(parts) < 2 {
				continue
			}

			index, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				continue
			}

			if index > lr.maxFiles {
				filesToDelete = append(filesToDelete, filepath.Join(dir, name))
			}
		}
	}

	for _, file := range filesToDelete {
		os.Remove(file)
	}

	return nil
}

func (lr *LogRotator) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 10, 5)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}