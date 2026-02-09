
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
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
    gzipExt      = ".gz"
)

type LogRotator struct {
    currentFile *os.File
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{basePath: basePath}
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
    rotatedPath := fmt.Sprintf("%s.%s%s", lr.basePath, timestamp, logExtension)
    if err := os.Rename(lr.basePath+logExtension, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressOldLogs(); err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath+logExtension, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

func (lr *LogRotator) compressOldLogs() error {
    pattern := lr.basePath + ".*" + logExtension
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    filesToCompress := matches[:len(matches)-maxBackups]
    for _, file := range filesToCompress {
        if err := lr.compressFile(file); err != nil {
            return err
        }
    }

    return nil
}

func (lr *LogRotator) compressFile(srcPath string) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstPath := srcPath + gzipExt
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

    if err := os.Remove(srcPath); err != nil {
        return err
    }

    return nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func extractTimestamp(filename string) (time.Time, error) {
    base := filepath.Base(filename)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return time.Time{}, fmt.Errorf("invalid filename format")
    }

    timestampStr := parts[1]
    return time.Parse("20060102_150405", timestampStr)
}

func getRotationNumber(filename string) (int, error) {
    base := filepath.Base(filename)
    parts := strings.Split(base, ".")
    if len(parts) < 2 {
        return 0, fmt.Errorf("invalid filename format")
    }

    numStr := strings.TrimSuffix(parts[len(parts)-1], gzipExt)
    return strconv.Atoi(numStr)
}

func main() {
    rotator, err := NewLogRotator("application")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingFile struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewRotatingFile(path string) (*RotatingFile, error) {
	rf := &RotatingFile{basePath: path}
	if err := rf.openCurrentFile(); err != nil {
		return nil, err
	}
	return rf, nil
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	if rf.currentSize+int64(len(p)) > maxFileSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.currentFile.Write(p)
	if err == nil {
		rf.currentSize += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) rotate() error {
	if rf.currentFile != nil {
		rf.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.gz", rf.basePath, timestamp)

	if err := compressFile(rf.basePath, backupPath); err != nil {
		return err
	}

	cleanupOldBackups(rf.basePath)

	return rf.openCurrentFile()
}

func compressFile(source, target string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	_, err = io.Copy(gz, src)
	return err
}

func cleanupOldBackups(basePath string) {
	pattern := basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		os.Remove(matches[i])
	}
}

func (rf *RotatingFile) openCurrentFile() error {
	file, err := os.OpenFile(rf.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rf.currentFile = file
	rf.currentSize = info.Size()
	return nil
}

func (rf *RotatingFile) Close() error {
	if rf.currentFile != nil {
		return rf.currentFile.Close()
	}
	return nil
}