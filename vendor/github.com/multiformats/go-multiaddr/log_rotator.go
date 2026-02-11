
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
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    filePath   string
    current    *os.File
    currentSize int64
}

func NewLogRotator(path string) (*LogRotator, error) {
    rotator := &LogRotator{filePath: path}
    if err := rotator.openCurrent(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (r *LogRotator) Write(p []byte) (int, error) {
    if r.currentSize+int64(len(p)) > maxFileSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.current.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *LogRotator) rotate() error {
    if err := r.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)
    if err := os.Rename(r.filePath, rotatedPath); err != nil {
        return err
    }

    if err := r.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := r.cleanupOldBackups(); err != nil {
        return err
    }

    return r.openCurrent()
}

func (r *LogRotator) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dstPath := path + ".gz"
    dst, err := os.Create(dstPath)
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(path); err != nil {
        return err
    }

    return nil
}

func (r *LogRotator) cleanupOldBackups() error {
    dir := filepath.Dir(r.filePath)
    base := filepath.Base(r.filePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, base+".") && strings.HasSuffix(name, ".gz") {
            backups = append(backups, name)
        }
    }

    if len(backups) <= maxBackups {
        return nil
    }

    for i := 0; i < len(backups)-maxBackups; i++ {
        path := filepath.Join(dir, backups[i])
        if err := os.Remove(path); err != nil {
            return err
        }
    }

    return nil
}

func (r *LogRotator) openCurrent() error {
    file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    r.current = file
    r.currentSize = info.Size()
    return nil
}

func (r *LogRotator) Close() error {
    if r.current != nil {
        return r.current.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: %s\n",
            time.Now().Format(time.RFC3339),
            i,
            strings.Repeat("X", 1024))

        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }

        if i%100 == 0 {
            fmt.Printf("Written %d entries\n", i)
        }
    }
}