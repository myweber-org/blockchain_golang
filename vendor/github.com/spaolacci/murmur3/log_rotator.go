
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
	currentSize int64
	fileCount   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
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

	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
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

	rl.fileCount++
	archiveName := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.fileCount, time.Now().Format("20060102T150405"))

	if err := rl.compressFile(rl.basePath, archiveName); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
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

type LogRotator struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compressOld bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (r *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(r.basePath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return err
    }

    file, err := os.OpenFile(r.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (r *LogRotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        err := r.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := r.currentFile.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *LogRotator) rotate() error {
    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

    err := os.Rename(r.basePath, backupPath)
    if err != nil {
        return err
    }

    err = r.openCurrentFile()
    if err != nil {
        return err
    }

    go r.manageBackups(backupPath)

    return nil
}

func (r *LogRotator) manageBackups(backupPath string) {
    if r.compressOld {
        compressedPath := backupPath + ".gz"
        err := compressFile(backupPath, compressedPath)
        if err == nil {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }

    backups, err := r.listBackups()
    if err != nil {
        return
    }

    if len(backups) > r.maxBackups {
        toRemove := backups[:len(backups)-r.maxBackups]
        for _, backup := range toRemove {
            os.Remove(backup)
        }
    }
}

func compressFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (r *LogRotator) listBackups() ([]string, error) {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            backups = append(backups, filepath.Join(dir, name))
        }
    }

    sortBackups(backups)
    return backups, nil
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(path string) int64 {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    if len(parts) < 2 {
        return 0
    }

    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }

    timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    return timestamp
}

func (r *LogRotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}