
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
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type LogRotator struct {
    mu          sync.Mutex
    currentFile *os.File
    currentSize int64
    basePath    string
    fileIndex   int
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

    oldPath := lr.currentFilePath()
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s.gz", oldPath, timestamp)

    if err := compressFile(oldPath, backupPath); err != nil {
        return err
    }

    os.Remove(oldPath)

    lr.fileIndex++
    if lr.fileIndex > maxBackupFiles {
        lr.fileIndex = 1
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    path := lr.currentFilePath()
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (lr *LogRotator) currentFilePath() string {
    return fmt.Sprintf("%s.%d.log", lr.basePath, lr.fileIndex)
}

func compressFile(source, target string) error {
    sourceFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    targetFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer targetFile.Close()

    gzWriter := gzip.NewWriter(targetFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    return err
}

func (lr *LogRotator) cleanupOldBackups() error {
    dir := filepath.Dir(lr.basePath)
    baseName := filepath.Base(lr.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName) && strings.HasSuffix(name, ".gz") {
            backups = append(backups, filepath.Join(dir, name))
        }
    }

    if len(backups) <= maxBackupFiles {
        return nil
    }

    sortBackups(backups)

    for i := maxBackupFiles; i < len(backups); i++ {
        os.Remove(backups[i])
    }

    return nil
}

func sortBackups(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            if extractTimestamp(files[i]) > extractTimestamp(files[j]) {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(filename string) string {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return ""
    }
    return parts[len(parts)-2]
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
    rotator, err := NewLogRotator("/var/log/myapp/application")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n",
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

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
	if maxSizeMB <= 0 {
		return nil, fmt.Errorf("maxSizeMB must be positive")
	}
	if maxFiles <= 0 {
		return nil, fmt.Errorf("maxFiles must be positive")
	}

	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
		maxFiles: maxFiles,
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
		rl.currentFile = nil
	}

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	if err := rl.compressFile(rl.basePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		if err := rl.cleanupOldFiles(); err != nil {
			return err
		}
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

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > rl.maxFiles {
		filesToDelete := matches[:len(matches)-rl.maxFiles]
		for _, f := range filesToDelete {
			if err := os.Remove(f); err != nil {
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
	logger, err := NewRotatingLogger("app.log", 10, 5)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}