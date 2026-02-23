package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentFile *os.File
    currentSize int64
    mu          sync.Mutex
}

func NewRotator(filePath string, maxSizeMB int, maxAgeHours int) (*Rotator, error) {
    absPath, err := filepath.Abs(filePath)
    if err != nil {
        return nil, err
    }

    dir := filepath.Dir(absPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    r := &Rotator{
        filePath: absPath,
        maxSize:  int64(maxSizeMB) * 1024 * 1024,
        maxAge:   time.Duration(maxAgeHours) * time.Hour,
    }

    if err := r.openCurrentFile(); err != nil {
        return nil, err
    }

    go r.cleanupOldFiles()
    return r, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.currentFile.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)
    if err := os.Rename(r.filePath, archivedPath); err != nil {
        return err
    }

    return r.openCurrentFile()
}

func (r *Rotator) openCurrentFile() error {
    file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.currentFile = file
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) cleanupOldFiles() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        r.mu.Lock()
        cutoff := time.Now().Add(-r.maxAge)
        dir := filepath.Dir(r.filePath)
        base := filepath.Base(r.filePath)

        entries, err := os.ReadDir(dir)
        if err != nil {
            r.mu.Unlock()
            continue
        }

        for _, entry := range entries {
            if entry.IsDir() {
                continue
            }

            name := entry.Name()
            if len(name) <= len(base)+1 || name[:len(base)] != base {
                continue
            }

            info, err := entry.Info()
            if err != nil {
                continue
            }

            if info.ModTime().Before(cutoff) {
                oldPath := filepath.Join(dir, name)
                os.Remove(oldPath)
            }
        }
        r.mu.Unlock()
    }
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("./logs/app.log", 10, 24)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        rotator.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    mu          sync.Mutex
    file        *os.File
    basePath    string
    maxSize     int64
    maxAge      time.Duration
    currentSize int64
    createdAt   time.Time
}

func NewRotator(basePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    r := &Rotator{
        basePath:  basePath,
        maxSize:   maxSize,
        maxAge:    maxAge,
        createdAt: time.Now(),
    }
    if err := r.openFile(); err != nil {
        return nil, err
    }
    return r, nil
}

func (r *Rotator) openFile() error {
    dir := filepath.Dir(r.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    f, err := os.OpenFile(r.basePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    r.file = f
    r.currentSize = info.Size()
    return nil
}

func (r *Rotator) rotate() error {
    if r.file != nil {
        r.file.Close()
    }
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    if err := os.Rename(r.basePath, backupPath); err != nil {
        return err
    }
    r.createdAt = time.Now()
    return r.openFile()
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize || time.Since(r.createdAt) > r.maxAge {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }
    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.file != nil {
        return r.file.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("./logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
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
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: name,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	if info, err := file.Stat(); err == nil {
		rl.currentSize = info.Size()
	}
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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
		if err := rl.compressCurrentFile(); err != nil {
			return err
		}
	}

	if err := rl.cleanOldBackups(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressCurrentFile() error {
	oldPath := rl.currentFile.Name()
	compressedPath := oldPath + ".gz"

	source, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, source); err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := filepath.Join(logDir, rl.baseName+"_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		toRemove := matches[:len(matches)-maxBackups]
		for _, file := range toRemove {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}