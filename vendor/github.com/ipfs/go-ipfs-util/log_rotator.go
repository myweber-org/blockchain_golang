
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
	backupCount = 5
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentSize int64
	file        *os.File
	basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == backupCount-1 {
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	currentBackup := rl.getBackupPath(0)
	if err := os.Rename(rl.basePath, currentBackup); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.basePath, index)
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingWriter struct {
	currentFile *os.File
	currentSize int64
	basePath    string
	backupCount int
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
	writer := &RotatingWriter{
		basePath: path,
	}
	if err := writer.openCurrentFile(); err != nil {
		return nil, err
	}
	return writer, nil
}

func (rw *RotatingWriter) Write(p []byte) (int, error) {
	if rw.currentSize+int64(len(p)) > maxFileSize {
		if err := rw.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rw.currentFile.Write(p)
	if err == nil {
		rw.currentSize += int64(n)
	}
	return n, err
}

func (rw *RotatingWriter) rotate() error {
	if rw.currentFile != nil {
		rw.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rw.basePath, timestamp)

	if err := os.Rename(rw.basePath, backupPath); err != nil {
		return err
	}

	rw.backupCount++
	if rw.backupCount > maxBackups {
		rw.cleanOldBackups()
	}

	return rw.openCurrentFile()
}

func (rw *RotatingWriter) openCurrentFile() error {
	file, err := os.OpenFile(rw.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rw.currentFile = file
	rw.currentSize = info.Size()
	return nil
}

func (rw *RotatingWriter) cleanOldBackups() {
	pattern := rw.basePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > maxBackups {
		for i := 0; i < len(matches)-maxBackups; i++ {
			os.Remove(matches[i])
		}
	}
}

func (rw *RotatingWriter) Close() error {
	if rw.currentFile != nil {
		return rw.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app.log")
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
			time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(logEntry))
		time.Sleep(10 * time.Millisecond)
	}
}