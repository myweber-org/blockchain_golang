
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingWriter struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingWriter(name string) (*RotatingWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	w := &RotatingWriter{baseName: basePath}

	if err := w.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := w.currentFile.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) rotateIfNeeded() error {
	if w.currentFile == nil || w.currentSize >= maxFileSize {
		return w.rotate()
	}
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.currentFile != nil {
		if err := w.currentFile.Close(); err != nil {
			return err
		}
		w.compressOldFile(w.currentFile.Name())
	}

	timestamp := time.Now().Format("20060102_150405")
	newPath := fmt.Sprintf("%s_%s.log", w.baseName, timestamp)

	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	w.currentFile = file
	if stat, err := file.Stat(); err == nil {
		w.currentSize = stat.Size()
	} else {
		w.currentSize = 0
	}

	w.cleanupOldFiles()
	return nil
}

func (w *RotatingWriter) compressOldFile(path string) {
	// Compression logic placeholder
	// In production, implement actual compression here
	fmt.Printf("Would compress: %s\n", path)
}

func (w *RotatingWriter) cleanupOldFiles() {
	pattern := fmt.Sprintf("%s_*.log", w.baseName)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		filesToRemove := matches[:len(matches)-maxBackups]
		for _, file := range filesToRemove {
			os.Remove(file)
		}
	}
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app")
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(logEntry))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation example completed")
}
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
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewLogRotator(baseName string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	lr := &LogRotator{
		baseName: baseName,
	}

	if err := lr.openCurrentFile(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", lr.baseName, timestamp))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = stat.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

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
	if lr.currentFile != nil {
		if err := lr.currentFile.Close(); err != nil {
			return err
		}

		oldPath := lr.currentFile.Name()
		compressedPath := oldPath + ".gz"

		if err := compressFile(oldPath, compressedPath); err != nil {
			return err
		}

		if err := os.Remove(oldPath); err != nil {
			return err
		}

		if err := lr.cleanupOldBackups(); err != nil {
			return err
		}
	}

	return lr.openCurrentFile()
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

func (lr *LogRotator) cleanupOldBackups() error {
	pattern := filepath.Join(logDir, lr.baseName+"_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackups {
		return nil
	}

	for i := 0; i < len(matches)-maxBackups; i++ {
		if err := os.Remove(matches[i]); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

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
		message := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
			time.Now().Format(time.RFC3339), i)
		rotator.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}