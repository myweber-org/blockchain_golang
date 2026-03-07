
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
    
    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    
    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)
    
    if err := os.Rename(lr.basePath, backupPath); err != nil {
        return err
    }
    
    if lr.compressOld {
        go lr.compressFile(backupPath)
    }
    
    if err := lr.cleanOldBackups(); err != nil {
        return err
    }
    
    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()
    
    destFile, err := os.Create(sourcePath + ".gz")
    if err != nil {
        return err
    }
    defer destFile.Close()
    
    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()
    
    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }
    
    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanOldBackups() error {
    pattern := lr.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }
    
    var backupFiles []string
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") || lr.isTimestampBackup(match) {
            backupFiles = append(backupFiles, match)
        }
    }
    
    if len(backupFiles) <= lr.maxBackups {
        return nil
    }
    
    filesToDelete := backupFiles[:len(backupFiles)-lr.maxBackups]
    for _, file := range filesToDelete {
        os.Remove(file)
    }
    
    return nil
}

func (lr *LogRotator) isTimestampBackup(path string) bool {
    suffix := strings.TrimPrefix(path, lr.basePath+".")
    if len(suffix) != 14 {
        return false
    }
    
    _, err := strconv.Atoi(suffix)
    return err == nil
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	baseName  string
	fileIndex int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrent(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	path := rl.currentPath()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) currentPath() string {
	return filepath.Join(logDir, rl.baseName+".log")
}

func (rl *RotatingLogger) backupPath(index int) string {
	timestamp := time.Now().Format("20060102-150405")
	return filepath.Join(logDir, fmt.Sprintf("%s.%d.%s.log.gz", rl.baseName, index, timestamp))
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	// Compress current log
	srcPath := rl.currentPath()
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	backupIndex := (rl.fileIndex % backupCount) + 1
	dstPath := rl.backupPath(backupIndex)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	// Remove original after compression
	if err := os.Remove(srcPath); err != nil {
		return err
	}

	// Open new log file
	if err := rl.openCurrent(); err != nil {
		return err
	}

	rl.fileIndex = backupIndex
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Redirect standard log output
	log.SetOutput(logger)

	// Generate log entries
	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, "Sample log message for testing rotation")
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}