
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

type RotatingWriter struct {
	filename   string
	current    *os.File
	size       int64
	mu         sync.Mutex
	maxSize    int64
	maxBackups int
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	w := &RotatingWriter{
		filename:   filename,
		maxSize:    maxFileSize,
		maxBackups: backupCount,
	}

	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openFile() error {
	dir := filepath.Dir(w.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.current = file
	w.size = stat.Size()
	return nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) >= w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.current.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if w.current != nil {
		if err := w.current.Close(); err != nil {
			return err
		}
	}

	for i := w.maxBackups - 1; i >= 0; i-- {
		oldName := w.backupName(i)
		newName := w.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(w.filename, w.backupName(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return w.openFile()
}

func (w *RotatingWriter) backupName(i int) string {
	if i == 0 {
		return w.filename + ".1"
	}
	return fmt.Sprintf("%s.%d", w.filename, i+1)
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.current != nil {
		return w.current.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log")
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := writer.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write failed: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

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
    maxSize       int64
    currentFile   *os.File
    currentSize   int64
    rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = info.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", l.basePath, timestamp)

    if err := os.Rename(l.basePath, rotatedPath); err != nil {
        return err
    }

    l.rotationCount++
    if err := l.compressOldLog(rotatedPath); err != nil {
        fmt.Printf("Failed to compress log: %v\n", err)
    }

    return l.openCurrentFile()
}

func (l *RotatingLogger) compressOldLog(path string) error {
    if l.rotationCount%5 != 0 {
        return nil
    }

    source, err := os.Open(path)
    if err != nil {
        return err
    }
    defer source.Close()

    compressedPath := path + ".gz"
    target, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer target.Close()

    gzWriter := gzip.NewWriter(target)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, source); err != nil {
        return err
    }

    os.Remove(path)
    return nil
}

func (l *RotatingLogger) cleanupOldFiles(maxAgeDays int) {
    cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
    pattern := l.basePath + ".*"

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    for _, match := range matches {
        info, err := os.Stat(match)
        if err != nil {
            continue
        }

        if info.ModTime().Before(cutoff) {
            os.Remove(match)
        }
    }
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func extractTimestamp(filename string) (time.Time, error) {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return time.Time{}, fmt.Errorf("invalid filename format")
    }

    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }

    return time.Parse("20060102_150405", timestampStr)
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()
        for range ticker.C {
            logger.cleanupOldFiles(30)
        }
    }()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation demonstration completed")
}