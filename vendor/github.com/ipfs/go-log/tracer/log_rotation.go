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