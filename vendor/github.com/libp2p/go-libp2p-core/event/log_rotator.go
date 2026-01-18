package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	file        *os.File
	filePath    string
	maxSize     int64
	currentSize int64
}

func NewRotatingLogger(filePath string, maxSizeMB int) (*RotatingLogger, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		file:        file,
		filePath:    absPath,
		maxSize:     int64(maxSizeMB) * 1024 * 1024,
		currentSize: info.Size(),
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
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
	if err := rl.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotatingLog struct {
	filePath   string
	maxSize    int64
	current    *os.File
	currentSize int64
}

func NewRotatingLog(path string, maxSizeMB int) (*RotatingLog, error) {
	rl := &RotatingLog{
		filePath: path,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	err := rl.openCurrent()
	return rl, err
}

func (rl *RotatingLog) openCurrent() error {
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.current = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLog) rotate() error {
	if rl.current != nil {
		rl.current.Close()
	}
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.filePath, timestamp)
	err := os.Rename(rl.filePath, archivePath+".tmp")
	if err != nil {
		return err
	}
	err = compressFile(archivePath+".tmp", archivePath)
	if err != nil {
		os.Rename(archivePath+".tmp", rl.filePath)
		return err
	}
	os.Remove(archivePath + ".tmp")
	return rl.openCurrent()
}

func compressFile(src, dst string) error {
	return nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		err := rl.rotate()
		if err != nil {
			return 0, err
		}
	}
	n, err := rl.current.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLog) Close() error {
	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	log, err := NewRotatingLog("app.log", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		log.Write([]byte(msg))
		time.Sleep(100 * time.Millisecond)
	}
}