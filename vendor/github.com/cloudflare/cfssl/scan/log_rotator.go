
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
	maxFileSize = 10 * 1024 * 1024
	backupCount = 5
)

type RotatingFile struct {
	mu         sync.Mutex
	filename   string
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingFile(filename string) (*RotatingFile, error) {
	rf := &RotatingFile{
		filename: filename,
		basePath: filepath.Dir(filename),
	}

	if err := rf.openCurrent(); err != nil {
		return nil, err
	}

	return rf, nil
}

func (rf *RotatingFile) openCurrent() error {
	info, err := os.Stat(rf.filename)
	if os.IsNotExist(err) {
		file, err := os.Create(rf.filename)
		if err != nil {
			return err
		}
		rf.file = file
		rf.size = 0
		return nil
	}
	if err != nil {
		return err
	}

	file, err := os.OpenFile(rf.filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rf.file = file
	rf.size = info.Size()
	return nil
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.size+int64(len(p)) > maxFileSize {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rf.file.Write(p)
	if err == nil {
		rf.size += int64(n)
	}
	return n, err
}

func (rf *RotatingFile) rotate() error {
	if err := rf.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s.gz", rf.filename, timestamp)

	if err := compressFile(rf.filename, backupName); err != nil {
		return err
	}

	if err := os.Remove(rf.filename); err != nil {
		return err
	}

	rf.cleanOldBackups()

	return rf.openCurrent()
}

func compressFile(source, target string) error {
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

func (rf *RotatingFile) cleanOldBackups() {
	pattern := rf.filename + ".*.gz"
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

func (rf *RotatingFile) Close() error {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.file.Close()
}

func main() {
	logFile, err := NewRotatingFile("application.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logFile.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}
}package main

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
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
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
		timestamp := time.Now().Format("20060102-150405")
		oldPath := filepath.Join(logDir, rl.baseName+".log")
		newPath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", rl.baseName, timestamp))
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
		if err := rl.compressFile(newPath); err != nil {
			return err
		}
		rl.cleanupOld()
	}
	return rl.openCurrent()
}

func (rl *RotatingLogger) openCurrent() error {
	path := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) compressFile(src string) error {
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

	gz := gzip.NewWriter(destFile)
	defer gz.Close()

	if _, err := io.Copy(gz, srcFile); err != nil {
		return err
	}
	os.Remove(src)
	return nil
}

func (rl *RotatingLogger) cleanupOld() {
	pattern := filepath.Join(logDir, rl.baseName+"-*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}
	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, f := range toDelete {
			os.Remove(f)
		}
	}
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
		msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(msg))
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
    "time"
)

const (
    maxLogSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
)

type RotatingLogger struct {
    filename    string
    currentSize int64
    file        *os.File
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filename: filename}
    if err := rl.openFile(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxLogSize {
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

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s.gz", rl.filename, timestamp)

    if err := compressFile(rl.filename, backupName); err != nil {
        return err
    }

    if err := cleanupOldBackups(rl.filename); err != nil {
        return err
    }

    return rl.openFile()
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

func cleanupOldBackups(baseFilename string) error {
    pattern := baseFilename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackupFiles {
        return nil
    }

    for i := 0; i < len(matches)-maxBackupFiles; i++ {
        if err := os.Remove(matches[i]); err != nil {
            return err
        }
    }
    return nil
}

func (rl *RotatingLogger) openFile() error {
    file, err := os.OpenFile(rl.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.file = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) Close() error {
    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}