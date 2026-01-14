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

type RotatingLogger struct {
	basePath      string
	maxSize       int64
	maxFiles      int
	currentSize   int64
	currentFile   *os.File
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
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
	}

	timestamp := time.Now().Unix()
	rotatedPath := fmt.Sprintf("%s.%d", rl.basePath, timestamp)
	if err := os.Rename(rl.basePath, rotatedPath); err != nil {
		return err
	}

	if err := rl.cleanupOldFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Cleanup error: %v\n", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	dir := filepath.Dir(rl.basePath)
	baseName := filepath.Base(rl.basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var rotatedFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, baseName+".") && entry.Type().IsRegular() {
			rotatedFiles = append(rotatedFiles, filepath.Join(dir, name))
		}
	}

	if len(rotatedFiles) <= rl.maxFiles {
		return nil
	}

	sort.Strings(rotatedFiles)
	filesToRemove := rotatedFiles[:len(rotatedFiles)-rl.maxFiles]

	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs/app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 1; i <= 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}