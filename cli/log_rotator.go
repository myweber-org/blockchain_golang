package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logExtension = ".log"
)

type RotatingLogger struct {
    currentFile *os.File
    basePath    string
    fileSize    int64
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    rl := &RotatingLogger{basePath: basePath}
    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    filename := rl.basePath + logExtension
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.fileSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.fileSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.fileSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s_%s%s", rl.basePath, timestamp, logExtension)

    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    if err := os.Rename(rl.basePath+logExtension, rotatedPath); err != nil {
        return err
    }

    if err := rl.openCurrentFile(); err != nil {
        return err
    }

    rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    files, err := os.ReadDir(dir)
    if err != nil {
        log.Printf("Failed to read directory: %v", err)
        return
    }

    var logFiles []string
    for _, file := range files {
        name := file.Name()
        if strings.HasPrefix(name, baseName) && strings.HasSuffix(name, logExtension) && name != filepath.Base(rl.basePath+logExtension) {
            logFiles = append(logFiles, filepath.Join(dir, name))
        }
    }

    if len(logFiles) <= maxBackups {
        return
    }

    for i := 0; i < len(logFiles)-maxBackups; i++ {
        if err := os.Remove(logFiles[i]); err != nil {
            log.Printf("Failed to remove old log file %s: %v", logFiles[i], err)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app_log")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    multiWriter := io.MultiWriter(os.Stdout, logger)
    log.SetOutput(multiWriter)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running smoothly", i)
        time.Sleep(10 * time.Millisecond)
    }
}package main

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
	maxFileSize = 10 * 1024 * 1024
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
		baseName: filepath.Join(logDir, name),
	}

	if err := rl.openCurrent(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	filename := rl.baseName + ".log"
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

	timestamp := time.Now().Format("20060102-150405")
	oldPath := rl.baseName + ".log"
	newPath := fmt.Sprintf("%s-%s.log", rl.baseName, timestamp)

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.compressFile(newPath); err != nil {
		log.Printf("Failed to compress %s: %v", newPath, err)
	}

	if err := rl.cleanupOld(); err != nil {
		log.Printf("Failed to cleanup old logs: %v", err)
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) compressFile(path string) error {
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

	return os.Remove(path)
}

func (rl *RotatingLogger) cleanupOld() error {
	pattern := rl.baseName + "-*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, file := range toDelete {
			if err := os.Remove(file); err != nil {
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}
}