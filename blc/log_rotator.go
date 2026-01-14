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

type RotatingLog struct {
	mu         sync.Mutex
	file       *os.File
	basePath   string
	maxSize    int64
	currentSize int64
	rotationCount int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
	rl := &RotatingLog{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLog) openCurrent() error {
	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLog) Write(p []byte) (int, error) {
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

func (rl *RotatingLog) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	rl.rotationCount++
	archiveName := fmt.Sprintf("%s.%d-%s.gz", rl.basePath, rl.rotationCount, time.Now().Format("20060102-150405"))

	if err := compressFile(rl.basePath, archiveName); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	return rl.openCurrent()
}

func compressFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
}

func (rl *RotatingLog) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	log, err := NewRotatingLog("app.log", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Application event processed successfully\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := log.Write([]byte(msg)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed. Check app.log and compressed archives.")
}