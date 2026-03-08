
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

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.currentFile = file
	rl.currentSize = 0

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
	}

	if err := rl.compressOldLogs(); err != nil {
		return err
	}

	if err := rl.cleanupOldBackups(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressOldLogs() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.log"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file) == ".gz" {
			continue
		}

		compressedFile := file + ".gz"
		if err := compressFile(file, compressedFile); err != nil {
			return err
		}

		os.Remove(file)
	}

	return nil
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

func (rl *RotatingLogger) cleanupOldBackups() error {
	files, err := filepath.Glob(filepath.Join(logDir, rl.baseName+"_*.gz"))
	if err != nil {
		return err
	}

	if len(files) <= maxBackups {
		return nil
	}

	for i := 0; i < len(files)-maxBackups; i++ {
		os.Remove(files[i])
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
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    lr.cleanupOldBackups()

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(source string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(source + ".gz")
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    return os.Remove(source)
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    backupFiles := make([]backupFile, 0, len(matches))
    for _, match := range matches {
        if ts, err := extractTimestamp(match); err == nil {
            backupFiles = append(backupFiles, backupFile{
                path:      match,
                timestamp: ts,
            })
        }
    }

    sortBackupsByTime(backupFiles)

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        os.Remove(backupFiles[i].path)
    }
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

type backupFile struct {
    path      string
    timestamp time.Time
}

func extractTimestamp(path string) (time.Time, error) {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return time.Time{}, fmt.Errorf("invalid backup file name")
    }

    timestampStr := parts[1]
    return time.Parse("20060102_150405", timestampStr)
}

func sortBackupsByTime(files []backupFile) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if files[i].timestamp.After(files[j].timestamp) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }

        if i%100 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
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
	mu         sync.Mutex
	file       *os.File
	size       int64
	baseName   string
	currentDay string
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, filename)
	rl := &RotatingLogger{
		baseName: basePath,
	}

	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDate {
		if rl.file != nil {
			rl.file.Close()
			if err := rl.compressOldLog(); err != nil {
				log.Printf("Failed to compress log: %v", err)
			}
			rl.cleanupOldBackups()
		}

		newPath := rl.baseName + "." + now.Format("20060102")
		if _, err := os.Stat(newPath); err == nil {
			for i := 1; ; i++ {
				newPath = fmt.Sprintf("%s.%s.%d", rl.baseName, now.Format("20060102"), i)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					break
				}
			}
		}

		file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		rl.file = file
		rl.size = 0
		rl.currentDay = currentDate
	}

	return nil
}

func (rl *RotatingLogger) compressOldLog() error {
	oldLogs, err := filepath.Glob(rl.baseName + ".*")
	if err != nil {
		return err
	}

	for _, logFile := range oldLogs {
		if filepath.Ext(logFile) == ".gz" {
			continue
		}

		compressedFile := logFile + ".gz"
		if _, err := os.Stat(compressedFile); err == nil {
			continue
		}

		src, err := os.Open(logFile)
		if err != nil {
			continue
		}

		dst, err := os.Create(compressedFile)
		if err != nil {
			src.Close()
			continue
		}

		gz := gzip.NewWriter(dst)
		if _, err := io.Copy(gz, src); err != nil {
			src.Close()
			gz.Close()
			dst.Close()
			continue
		}

		src.Close()
		gz.Close()
		dst.Close()
		os.Remove(logFile)
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
	backups, err := filepath.Glob(rl.baseName + ".*.gz")
	if err != nil {
		return
	}

	if len(backups) > backupCount {
		sortBackups(backups)
		for i := backupCount; i < len(backups); i++ {
			os.Remove(backups[i])
		}
	}
}

func sortBackups(files []string) {
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			infoI, _ := os.Stat(files[i])
			infoJ, _ := os.Stat(files[j])
			if infoI.ModTime().Before(infoJ.ModTime()) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
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
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: Processing request at %v", i, time.Now())
		time.Sleep(100 * time.Millisecond)
	}
}
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
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
		file:        file,
		basePath:    basePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		backupCount: backupCount,
		compressOld: compressOld,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	ext := filepath.Ext(rl.basePath)
	baseName := strings.TrimSuffix(rl.basePath, ext)
	timestamp := time.Now().Format("20060102_150405")

	for i := rl.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		
		if i == 0 {
			oldPath = rl.basePath
		} else {
			oldPath = fmt.Sprintf("%s.%d%s", baseName, i, ext)
		}
		
		if rl.compressOld && i == rl.backupCount-1 {
			newPath = fmt.Sprintf("%s.%d.%s.gz", baseName, i, timestamp)
		} else {
			newPath = fmt.Sprintf("%s.%d%s", baseName, i+1, ext)
		}

		if _, err := os.Stat(oldPath); err == nil {
			if rl.compressOld && i == rl.backupCount-1 {
				if err := compressFile(oldPath, newPath); err != nil {
					log.Printf("Failed to compress %s: %v", oldPath, err)
					os.Rename(oldPath, newPath)
				}
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.currentSize = 0
	return nil
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

	gzWriter := NewGzipWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	if err != nil {
		os.Remove(dst)
		return err
	}

	os.Remove(src)
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

type GzipWriter struct {
	io.WriteCloser
}

func NewGzipWriter(w io.Writer) *GzipWriter {
	return &GzipWriter{}
}

func (gw *GzipWriter) Close() error {
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d: Application event occurred", i)
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
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewLogRotator(baseName string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	lr := &LogRotator{
		baseName: baseName,
	}

	if err := lr.openCurrentFile(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", lr.baseName, timestamp))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	if info, err := file.Stat(); err == nil {
		lr.currentSize = info.Size()
	}
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > maxFileSize {
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

	oldLogs, err := filepath.Glob(filepath.Join(logDir, lr.baseName+"_*.log"))
	if err != nil {
		return err
	}

	if len(oldLogs) >= maxBackups {
		for i := 0; i <= len(oldLogs)-maxBackups; i++ {
			os.Remove(oldLogs[i])
		}
	}

	for _, logFile := range oldLogs {
		if err := compressLog(logFile); err != nil {
			fmt.Printf("Failed to compress %s: %v\n", logFile, err)
		}
	}

	return lr.openCurrentFile()
}

func compressLog(filename string) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(filename)
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
	rotator, err := NewLogRotator("app")
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
			time.Now().Format(time.RFC3339), i)
		rotator.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}