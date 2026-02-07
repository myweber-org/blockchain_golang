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
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rl := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.currentFile = f
	if info, err := f.Stat(); err == nil {
		rl.currentSize = info.Size()
	}
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

	rl.rotationCount++
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.%d.gz", rl.filePath, timestamp, rl.rotationCount)

	if err := rl.compressFile(rl.filePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(rl.filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	rl.currentSize = 0
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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
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

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
    gzipExt      = ".gz"
)

type RotatingLogger struct {
    mu        sync.Mutex
    file      *os.File
    size      int64
    basePath  string
    currentID int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath:  strings.TrimSuffix(path, logExtension),
    }

    if err := rl.openOrCreate(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.file != nil {
        rl.file.Close()
    }

    rl.currentID++

    if rl.currentID > maxBackups {
        rl.cleanOldLogs()
    }

    return rl.openOrCreate()
}

func (rl *RotatingLogger) openOrCreate() error {
    filename := rl.basePath + logExtension
    if rl.currentID > 0 {
        filename = fmt.Sprintf("%s.%d%s", rl.basePath, rl.currentID, logExtension)
    }

    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

    if rl.currentID > 0 {
        go rl.compressPreviousLog()
    }

    return nil
}

func (rl *RotatingLogger) compressPreviousLog() {
    prevID := rl.currentID - 1
    if prevID <= 0 {
        return
    }

    srcName := fmt.Sprintf("%s.%d%s", rl.basePath, prevID, logExtension)
    dstName := srcName + gzipExt

    srcFile, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return
    }

    os.Remove(srcName)
}

func (rl *RotatingLogger) cleanOldLogs() {
    files, err := filepath.Glob(rl.basePath + ".*" + logExtension + "*")
    if err != nil {
        return
    }

    var logFiles []string
    for _, f := range files {
        if strings.HasSuffix(f, gzipExt) || strings.HasSuffix(f, logExtension) {
            logFiles = append(logFiles, f)
        }
    }

    if len(logFiles) <= maxBackups {
        return
    }

    sortFiles(logFiles)
    for i := 0; i < len(logFiles)-maxBackups; i++ {
        os.Remove(logFiles[i])
    }
}

func sortFiles(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractID(files[i]) > extractID(files[j]) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractID(filename string) int {
    base := filepath.Base(filename)
    base = strings.TrimSuffix(base, logExtension)
    base = strings.TrimSuffix(base, gzipExt)

    parts := strings.Split(base, ".")
    if len(parts) < 2 {
        return 0
    }

    id, _ := strconv.Atoi(parts[len(parts)-1])
    return id
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("application.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}