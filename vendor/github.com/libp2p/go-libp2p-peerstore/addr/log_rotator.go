package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	maxFiles    int
	currentSize int64
	file        *os.File
}

func NewLogRotator(filePath string, maxSizeMB int, maxFiles int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	dir := filepath.Dir(lr.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to get file info: %w", err)
	}

	lr.file = file
	lr.currentSize = info.Size()

	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	n, err := lr.file.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedFile := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, rotatedFile); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if err := lr.compressFile(rotatedFile); err != nil {
		log.Printf("Warning: failed to compress %s: %v", rotatedFile, err)
	}

	if err := lr.cleanupOldFiles(); err != nil {
		log.Printf("Warning: failed to cleanup old files: %v", err)
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(filePath string) error {
	source, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer source.Close()

	compressedPath := filePath + ".gz"
	dest, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gzWriter := gzip.NewWriter(dest)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, source); err != nil {
		os.Remove(compressedPath)
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func (lr *LogRotator) cleanupOldFiles() error {
	dir := filepath.Dir(lr.filePath)
	baseName := filepath.Base(lr.filePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var rotatedFiles []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), baseName+".") {
			rotatedFiles = append(rotatedFiles, entry.Name())
		}
	}

	if len(rotatedFiles) <= lr.maxFiles {
		return nil
	}

	sort.Slice(rotatedFiles, func(i, j int) bool {
		return rotatedFiles[i] > rotatedFiles[j]
	})

	for i := lr.maxFiles; i < len(rotatedFiles); i++ {
		fileToRemove := filepath.Join(dir, rotatedFiles[i])
		if err := os.Remove(fileToRemove); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) Close() error {
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(rotator)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: Application is running normally", i)
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

type RotatingLogger struct {
    mu           sync.Mutex
    currentFile  *os.File
    basePath     string
    maxSize      int64
    currentSize  int64
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
        currentFile:  file,
        basePath:     basePath,
        maxSize:      maxSize,
        currentSize:  info.Size(),
        rotationCount: 0,
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
    
    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    rl.currentFile = file
    rl.currentSize = 0
    rl.rotationCount++
    
    return nil
}

func (rl *RotatingLogger) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()
    
    compressedPath := sourcePath + ".gz"
    destFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer destFile.Close()
    
    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()
    
    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }
    
    if err := os.Remove(sourcePath); err != nil {
        return err
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

func (rl *RotatingLogger) CleanupOldLogs(maxAge time.Duration) error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    
    cutoff := time.Now().Add(-maxAge)
    
    for _, file := range matches {
        info, err := os.Stat(file)
        if err != nil {
            continue
        }
        
        if info.ModTime().Before(cutoff) {
            os.Remove(file)
        }
    }
    
    return nil
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
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentSize int64
    createdAt   time.Time
}

func NewRotator(filePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    dir := filepath.Dir(filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create directory: %w", err)
    }

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to stat log file: %w", err)
    }

    return &Rotator{
        file:        file,
        filePath:    filePath,
        maxSize:     maxSize,
        maxAge:      maxAge,
        currentSize: info.Size(),
        createdAt:   time.Now(),
    }, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.shouldRotate() {
        if err := r.rotate(); err != nil {
            return 0, fmt.Errorf("rotate failed: %w", err)
        }
    }

    n, err := r.file.Write(p)
    if err != nil {
        return n, fmt.Errorf("write failed: %w", err)
    }
    r.currentSize += int64(n)
    return n, nil
}

func (r *Rotator) shouldRotate() bool {
    if r.maxSize > 0 && r.currentSize >= r.maxSize {
        return true
    }
    if r.maxAge > 0 && time.Since(r.createdAt) >= r.maxAge {
        return true
    }
    return false
}

func (r *Rotator) rotate() error {
    if err := r.file.Close(); err != nil {
        return fmt.Errorf("failed to close file: %w", err)
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return fmt.Errorf("failed to rename file: %w", err)
    }

    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    r.file = file
    r.currentSize = 0
    r.createdAt = time.Now()
    return nil
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.file.Close()
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Printf("Failed to create rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation example completed")
}