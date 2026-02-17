
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

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	currentSize int64
	fileCount   int
	maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
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

	f, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rl.currentFile = f
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.currentFile.Close()

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	source, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gzWriter := gzip.NewWriter(dest)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, source); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		rl.cleanupOldFiles()
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
	files, err := filepath.Glob(rl.basePath + ".*.gz")
	if err != nil {
		return
	}

	if len(files) > rl.maxFiles {
		filesToDelete := files[:len(files)-rl.maxFiles]
		for _, f := range filesToDelete {
			os.Remove(f)
		}
	}
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d: Application event occurred at %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func rotateLog(logPath string, maxBackups int) error {
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s", logPath, timestamp)

	err = os.Rename(logPath, archivePath)
	if err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	matches, err := filepath.Glob(logPath + ".*")
	if err != nil {
		return fmt.Errorf("failed to list backup files: %w", err)
	}

	if len(matches) > maxBackups {
		oldestFiles := matches[:len(matches)-maxBackups]
		for _, oldFile := range oldestFiles {
			err = os.Remove(oldFile)
			if err != nil {
				return fmt.Errorf("failed to remove old backup %s: %w", oldFile, err)
			}
		}
	}

	newFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	newFile.Close()

	fmt.Printf("Log rotated: %s -> %s\n", logPath, archivePath)
	return nil
}

func main() {
	logFile := "application.log"
	err := rotateLog(logFile, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rotating log: %v\n", err)
		os.Exit(1)
	}
}