
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
	fileSize    int64
	fileCount   int
	maxFiles    int
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
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
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.fileSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.fileSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	if err := rl.compressFile(rl.basePath, archivePath); err != nil {
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

func (rl *RotatingLogger) compressFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	_, err = io.Copy(gz, source)
	return err
}

func (rl *RotatingLogger) cleanupOldFiles() {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxFiles {
		filesToRemove := len(matches) - rl.maxFiles
		for i := 0; i < filesToRemove; i++ {
			os.Remove(matches[i])
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 1024*1024, 10)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogRotator struct {
	CurrentLogPath string
	MaxSize        int64
	ArchiveDir     string
}

func NewLogRotator(logPath string, maxSize int64, archiveDir string) *LogRotator {
	return &LogRotator{
		CurrentLogPath: logPath,
		MaxSize:        maxSize,
		ArchiveDir:     archiveDir,
	}
}

func (lr *LogRotator) CheckAndRotate() error {
	info, err := os.Stat(lr.CurrentLogPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if info.Size() < lr.MaxSize {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveName := filepath.Base(lr.CurrentLogPath) + "." + timestamp
	archivePath := filepath.Join(lr.ArchiveDir, archiveName)

	if err := os.Rename(lr.CurrentLogPath, archivePath); err != nil {
		return fmt.Errorf("failed to archive log file: %w", err)
	}

	newFile, err := os.Create(lr.CurrentLogPath)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	newFile.Close()

	fmt.Printf("Log rotated: %s -> %s\n", lr.CurrentLogPath, archivePath)
	return nil
}

func main() {
	rotator := NewLogRotator("/var/log/myapp/app.log", 10*1024*1024, "/var/log/myapp/archive")
	if err := rotator.CheckAndRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error rotating logs: %v\n", err)
		os.Exit(1)
	}
}