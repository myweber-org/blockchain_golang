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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	sequence    int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	oldPath := filepath.Join(logDir, rl.baseName+".log")
	newPath := filepath.Join(logDir, rl.baseName+"."+time.Now().Format("20060102150405")+".log")

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.compressOldLogs(); err != nil {
		log.Printf("Failed to compress logs: %v", err)
	}

	if err := rl.cleanupOldLogs(); err != nil {
		log.Printf("Failed to cleanup old logs: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressOldLogs() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+".*.log"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".gz") {
			continue
		}

		if err := compressFile(file); err != nil {
			return err
		}
	}

	return nil
}

func compressFile(src string) error {
	dest := src + ".gz"
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Simple copy for demonstration (in real implementation use gzip.Writer)
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		os.Remove(dest)
		return err
	}

	if err := os.Remove(src); err != nil {
		os.Remove(dest)
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldLogs() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+".*.gz"))
	if err != nil {
		return err
	}

	if len(files) <= maxBackups {
		return nil
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(files)-maxBackups; i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}

	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
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

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compress bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compress,
    }
    
    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (r *LogRotator) openCurrentFile() error {
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
    rotatedPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)
    
    err := os.Rename(r.basePath, rotatedPath)
    if err != nil {
        return err
    }
    
    if r.compressOld {
        go r.compressFile(rotatedPath)
    }
    
    err = r.openCurrentFile()
    if err != nil {
        return err
    }
    
    r.cleanupOldBackups()
    
    return nil
}

func (r *LogRotator) compressFile(path string) {
    srcFile, err := os.Open(path)
    if err != nil {
        return
    }
    defer srcFile.Close()
    
    destPath := path + ".gz"
    destFile, err := os.Create(destPath)
    if err != nil {
        return
    }
    defer destFile.Close()
    
    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()
    
    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return
    }
    
    os.Remove(path)
}

func (r *LogRotator) cleanupOldBackups() {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)
    
    files, err := os.ReadDir(dir)
    if err != nil {
        return
    }
    
    var backupFiles []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, baseName+".") {
            backupFiles = append(backupFiles, name)
        }
    }
    
    if len(backupFiles) <= r.maxBackups {
        return
    }
    
    sortBackupFiles(backupFiles)
    
    for i := 0; i < len(backupFiles)-r.maxBackups; i++ {
        os.Remove(filepath.Join(dir, backupFiles[i]))
    }
}

func sortBackupFiles(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractTimestamp(files[i]) > extractTimestamp(files[j]) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    parts := strings.Split(filename, ".")
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

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()
    
    for i := 0; i < 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}