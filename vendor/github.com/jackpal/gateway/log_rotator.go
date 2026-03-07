
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
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu         sync.Mutex
	current    *os.File
	baseName   string
	fileSize   int64
	sequence   int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logger := &RotatingLogger{
		baseName: baseName,
		sequence: 0,
	}

	if err := logger.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.current.Write(p)
	if err != nil {
		return n, err
	}

	rl.fileSize += int64(n)

	if err := rl.rotateIfNeeded(); err != nil {
		return n, err
	}

	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.current == nil || rl.fileSize >= maxFileSize {
		return rl.rotate()
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	if rl.current != nil {
		rl.current.Close()
		if err := rl.compressOldLog(); err != nil {
			fmt.Printf("Failed to compress log: %v\n", err)
		}
	}

	rl.sequence++
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.current = file
	rl.fileSize = 0

	return nil
}

func (rl *RotatingLogger) compressOldLog() error {
	if rl.sequence <= 1 {
		return nil
	}

	oldName := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.sequence-1))
	compressedName := oldName + ".gz"

	oldFile, err := os.Open(oldName)
	if err != nil {
		return err
	}
	defer oldFile.Close()

	compressedFile, err := os.Create(compressedName)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	gzWriter := gzip.NewWriter(compressedFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		return err
	}

	os.Remove(oldName)
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu            sync.Mutex
    basePath      string
    currentFile   *os.File
    maxSize       int64
    currentSize   int64
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        basePath:    basePath,
        currentFile: file,
        maxSize:     maxSize,
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0
    rl.rotationCount++

    go rl.compressOldLog(rotatedPath)
    go rl.cleanupOldLogs()

    return nil
}

func (rl *RotatingLogger) compressOldLog(path string) {
    src, err := os.Open(path)
    if err != nil {
        return
    }
    defer src.Close()

    dstPath := path + ".gz"
    dst, err := os.Create(dstPath)
    if err != nil {
        return
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return
    }

    os.Remove(path)
}

func (rl *RotatingLogger) cleanupOldLogs() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var logs []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            logs = append(logs, filepath.Join(dir, name))
        }
    }

    if len(logs) > 10 {
        logs = logs[:len(logs)-10]
        for _, log := range logs {
            os.Remove(log)
        }
    }
}

func (rl *RotatingLogger) parseRotationNumber(filename string) int {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return 0
    }
    num, _ := strconv.Atoi(parts[len(parts)-2])
    return num
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
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}