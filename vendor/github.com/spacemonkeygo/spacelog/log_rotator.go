package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupLogs = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    filePath    string
    currentSize int64
}

func NewRotatingLogger(baseDir string) (*RotatingLogger, error) {
    filePath := filepath.Join(baseDir, logFileName)
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentFile: file,
        filePath:    filePath,
        currentSize: info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.currentSize+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.currentFile.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0

    go rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    dir := filepath.Dir(rl.filePath)
    pattern := filepath.Join(dir, logFileName+".*")

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackupLogs {
        filesToDelete := matches[:len(matches)-maxBackupLogs]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: This is a test log message for rotation testing", i)
        time.Sleep(10 * time.Millisecond)
    }
}package main

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

type RotatingLogger struct {
	basePath   string
	maxSize    int64
	maxBackups int
	current    *os.File
	size       int64
}

func NewRotatingLogger(path string, maxSizeMB int, maxBackups int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rl := &RotatingLogger{
		basePath:   path,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}

	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	f, err := os.OpenFile(rl.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rl.current = f
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.size+int64(len(p)) >= rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.current.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.current.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")
	backupPath := rl.basePath + "." + timestamp

	if err := os.Rename(rl.basePath, backupPath); err != nil {
		return err
	}

	if err := rl.openCurrent(); err != nil {
		return err
	}

	go rl.cleanupOldBackups()
	go rl.compressBackup(backupPath)

	return nil
}

func (rl *RotatingLogger) compressBackup(path string) {
	compressedPath := path + ".gz"
	log.Printf("Compressing %s to %s", path, compressedPath)
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := rl.basePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= rl.maxBackups {
		return
	}

	backups := make([]string, 0, len(matches))
	for _, match := range matches {
		if strings.HasSuffix(match, ".gz") {
			backups = append(backups, match)
		}
	}

	sortBackupsByTime(backups)

	for i := 0; i < len(backups)-rl.maxBackups; i++ {
		os.Remove(backups[i])
	}
}

func sortBackupsByTime(backups []string) {
	for i := 0; i < len(backups); i++ {
		for j := i + 1; j < len(backups); j++ {
			if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}
}

func extractTimestamp(path string) int64 {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return 0
	}

	timestampStr := parts[len(parts)-2]
	if strings.HasSuffix(timestampStr, "gz") && len(parts) >= 3 {
		timestampStr = parts[len(parts)-3]
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return 0
	}
	return timestamp
}

func (rl *RotatingLogger) Close() error {
	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
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

type RotatingLog struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
    rotationSeq int
}

func NewRotatingLog(path string, maxSizeMB int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLog{
        filePath:    path,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
        rotationSeq: 0,
    }, nil
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
    if rl.file != nil {
        rl.file.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s.%d", rl.filePath, timestamp, rl.rotationSeq)
    rl.rotationSeq++

    if err := os.Rename(rl.filePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    info, err := file.Stat()
    if err != nil {
        return err
    }

    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLog) compressFile(src string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dst := src + ".gz"
    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gz := gzip.NewWriter(dstFile)
    defer gz.Close()

    if _, err := io.Copy(gz, srcFile); err != nil {
        return err
    }

    os.Remove(src)
    return nil
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}