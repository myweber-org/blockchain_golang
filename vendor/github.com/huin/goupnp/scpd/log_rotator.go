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

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentFile *os.File
    currentSize int64
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = stat.Size()
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

    if err := l.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := l.cleanupOldFiles(); err != nil {
        return err
    }

    return l.openCurrentFile()
}

func (l *RotatingLogger) compressFile(srcPath string) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstPath := srcPath + ".gz"
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

    os.Remove(srcPath)
    return nil
}

func (l *RotatingLogger) cleanupOldFiles() error {
    pattern := l.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= l.maxFiles {
        return nil
    }

    var timestamps []time.Time
    fileMap := make(map[time.Time]string)

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        timestampStr := parts[len(parts)-2]
        t, err := time.Parse("20060102_150405", timestampStr)
        if err != nil {
            continue
        }
        timestamps = append(timestamps, t)
        fileMap[t] = match
    }

    for i := 0; i < len(timestamps)-l.maxFiles; i++ {
        oldestTime := timestamps[i]
        if filePath, exists := fileMap[oldestTime]; exists {
            os.Remove(filePath)
        }
    }

    return nil
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(message))
        time.Sleep(100 * time.Millisecond)
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
    "sync"
)

type LogRotator struct {
    basePath     string
    maxSize      int64
    maxBackups   int
    currentSize  int64
    file         *os.File
    mu           sync.Mutex
}

func NewLogRotator(basePath string, maxSize int64, maxBackups int) (*LogRotator, error) {
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    for i := lr.maxBackups - 1; i >= 0; i-- {
        oldPath := lr.backupPath(i)
        if _, err := os.Stat(oldPath); os.IsNotExist(err) {
            continue
        }

        newPath := lr.backupPath(i + 1)
        if i == lr.maxBackups-1 {
            if err := os.Remove(newPath); err != nil && !os.IsNotExist(err) {
                return err
            }
            continue
        }

        if err := os.Rename(oldPath, newPath); err != nil {
            return err
        }
    }

    if err := os.Rename(lr.basePath, lr.backupPath(0)); err != nil && !os.IsNotExist(err) {
        return err
    }

    file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0
    return nil
}

func (lr *LogRotator) backupPath(index int) string {
    if index == 0 {
        return lr.basePath + ".1"
    }
    return lr.basePath + "." + strconv.Itoa(index+1) + ".gz"
}

func (lr *LogRotator) compressOldLogs() error {
    for i := 0; i < lr.maxBackups; i++ {
        path := lr.backupPath(i)
        if i == 0 {
            if err := lr.compressFile(path, path+".gz"); err != nil && !os.IsNotExist(err) {
                return err
            }
            if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
                return err
            }
        }
    }
    return nil
}

func (lr *LogRotator) compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator("app.log", 1024*1024, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logLine := fmt.Sprintf("Log entry %d: Application event recorded\n", i)
        if _, err := rotator.Write([]byte(logLine)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
    }

    if err := rotator.compressOldLogs(); err != nil {
        fmt.Printf("Compression error: %v\n", err)
    }
}
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type RotatingFile struct {
	MaxSize    int64
	BackupCount int
	Compress   bool
	basePath   string
	current   *os.File
	size      int64
}

func NewRotatingFile(path string, maxSize int64, backupCount int, compress bool) (*RotatingFile, error) {
	rf := &RotatingFile{
		MaxSize:     maxSize,
		BackupCount: backupCount,
		Compress:    compress,
		basePath:    path,
	}

	if err := rf.openCurrent(); err != nil {
		return nil, err
	}

	return rf, nil
}

func (rf *RotatingFile) openCurrent() error {
	if rf.current != nil {
		rf.current.Close()
	}

	f, err := os.OpenFile(rf.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rf.current = f
	rf.size = stat.Size()
	return nil
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	if rf.size+int64(len(p)) >= rf.MaxSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.current.Write(p)
	if err == nil {
		rf.size += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) rotate() error {
	if err := rf.current.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rf.basePath, timestamp)

	if err := os.Rename(rf.basePath, backupPath); err != nil {
		return err
	}

	if rf.Compress {
		go rf.compressBackup(backupPath)
	}

	if err := rf.cleanOldBackups(); err != nil {
		log.Printf("Failed to clean old backups: %v", err)
	}

	return rf.openCurrent()
}

func (rf *RotatingFile) compressBackup(path string) {
	// Compression implementation would go here
	// For now just log the operation
	log.Printf("Compressing backup: %s", path)
}

func (rf *RotatingFile) cleanOldBackups() error {
	pattern := rf.basePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= rf.BackupCount {
		return nil
	}

	toRemove := matches[:len(matches)-rf.BackupCount]
	for _, path := range toRemove {
		if err := os.Remove(path); err != nil {
			return err
		}
		log.Printf("Removed old backup: %s", path)
	}

	return nil
}

func (rf *RotatingFile) Close() error {
	if rf.current != nil {
		return rf.current.Close()
	}
	return nil
}

func main() {
	logFile, err := NewRotatingFile("app.log", 1024*1024, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d: Application is running normally", i)
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
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
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

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	if err := rl.compressFile(rl.basePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	rl.fileCount++
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
	logger, err := NewRotatingLogger("app.log", 10)
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