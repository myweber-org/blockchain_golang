
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingWriter struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	fileCount   int
	maxFiles    int
}

func NewRotatingWriter(basePath string, maxSize int64, maxFiles int) (*RotatingWriter, error) {
	dir := filepath.Dir(basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingWriter{
		currentFile: file,
		filePath:    basePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		maxFiles:    maxFiles,
	}, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.currentFile.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if err := w.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s", w.filePath, timestamp)

	if err := os.Rename(w.filePath, archivePath); err != nil {
		return err
	}

	file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.currentFile = file
	w.currentSize = 0
	w.fileCount++

	if w.fileCount > w.maxFiles {
		w.cleanupOldFiles()
	}

	return nil
}

func (w *RotatingWriter) cleanupOldFiles() {
	dir := filepath.Dir(w.filePath)
	baseName := filepath.Base(w.filePath)

	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var archiveFiles []string
	for _, file := range files {
		name := file.Name()
		if len(name) > len(baseName) && name[:len(baseName)] == baseName && name[len(baseName)] == '.' {
			archiveFiles = append(archiveFiles, filepath.Join(dir, name))
		}
	}

	if len(archiveFiles) <= w.maxFiles {
		return
	}

	for i := 0; i < len(archiveFiles)-w.maxFiles; i++ {
		os.Remove(archiveFiles[i])
	}
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.currentFile.Close()
}

func main() {
	writer, err := NewRotatingWriter("/var/log/myapp/app.log", 1024*1024, 10)
	if err != nil {
		fmt.Printf("Failed to create rotating writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		writer.Write([]byte(message))
		time.Sleep(100 * time.Millisecond)
	}
}