package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxFiles    = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu       sync.Mutex
	file     *os.File
	baseName string
	counter  int
}

func NewRotatingLogger() (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	baseName := filepath.Join(logDir, "app")
	rl := &RotatingLogger{baseName: baseName}
	if err := rl.rotate(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	info, err := rl.file.Stat()
	if err != nil {
		return 0, err
	}

	if info.Size()+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	return rl.file.Write(p)
}

func (rl *RotatingLogger) rotate() error {
	if rl.file != nil {
		rl.file.Close()
	}

	rl.counter = (rl.counter + 1) % maxFiles
	filename := rl.baseName + "_" + strconv.Itoa(rl.counter) + ".log"

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	rl.file = file

	go rl.cleanupOldFiles()
	return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
	files, err := filepath.Glob(rl.baseName + "_*.log")
	if err != nil {
		return
	}

	if len(files) > maxFiles {
		for _, file := range files[:len(files)-maxFiles] {
			os.Remove(file)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d at %v", i, time.Now())
		time.Sleep(100 * time.Millisecond)
	}
}
package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logFileName = "app.log"
)

type RotatingLogger struct {
	currentSize int64
	backupCount int
}

func NewRotatingLogger() *RotatingLogger {
	return &RotatingLogger{}
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.checkRotation()
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err = file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) checkRotation() {
	info, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		rl.currentSize = 0
		return
	}
	if err != nil {
		log.Printf("Error checking log file: %v", err)
		return
	}

	rl.currentSize = info.Size()
	if rl.currentSize >= maxFileSize {
		rl.rotate()
	}
}

func (rl *RotatingLogger) rotate() {
	timestamp := time.Now().Format("20060102_150405")
	backupName := logFileName + "." + timestamp

	err := os.Rename(logFileName, backupName)
	if err != nil {
		log.Printf("Failed to rotate log: %v", err)
		return
	}

	rl.currentSize = 0
	rl.cleanupOldBackups()
	rl.backupCount++
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := logFileName + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackups {
		return
	}

	backups := make([]backupInfo, 0, len(matches))
	for _, match := range matches {
		parts := strings.Split(match, ".")
		if len(parts) < 2 {
			continue
		}
		timestamp := parts[len(parts)-1]
		t, err := time.Parse("20060102_150405", timestamp)
		if err != nil {
			continue
		}
		backups = append(backups, backupInfo{path: match, time: t})
	}

	if len(backups) <= maxBackups {
		return
	}

	for i := 0; i < len(backups)-maxBackups; i++ {
		os.Remove(backups[i].path)
	}
}

type backupInfo struct {
	path string
	time time.Time
}

func main() {
	logger := NewRotatingLogger()
	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strconv.FormatInt(time.Now().UnixNano(), 10))
		time.Sleep(10 * time.Millisecond)
	}
}