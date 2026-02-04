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
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	filePath    string
	mu          sync.Mutex
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: basePath,
	}
	if err := rl.openCurrentFile(); err != nil {
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
	}

	timestamp := time.Now().Format("20060102-150405")
	oldPath := rl.filePath + "." + timestamp
	if err := os.Rename(rl.filePath, oldPath); err != nil {
		return err
	}

	if err := rl.compressFile(oldPath); err != nil {
		return err
	}

	if err := rl.cleanupOldFiles(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(src + ".gz")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	if err != nil {
		return err
	}

	return os.Remove(src)
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	files, err := filepath.Glob(rl.filePath + ".*.gz")
	if err != nil {
		return err
	}

	if len(files) > maxBackups {
		filesToRemove := files[:len(files)-maxBackups]
		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
)

type RotatingFile struct {
    mu         sync.Mutex
    file       *os.File
    basePath   string
    maxSize    int64
    currentSize int64
    fileIndex  int
}

func NewRotatingFile(basePath string, maxSize int64) (*RotatingFile, error) {
    rf := &RotatingFile{
        basePath:  basePath,
        maxSize:   maxSize,
        fileIndex: 0,
    }
    
    if err := rf.openFile(); err != nil {
        return nil, err
    }
    
    return rf, nil
}

func (rf *RotatingFile) openFile() error {
    filename := rf.basePath
    if rf.fileIndex > 0 {
        filename = rf.basePath + "." + strconv.Itoa(rf.fileIndex)
    }
    
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    rf.file = file
    rf.currentSize = info.Size()
    return nil
}

func (rf *RotatingFile) rotate() error {
    rf.file.Close()
    rf.fileIndex++
    
    for i := rf.fileIndex; i > 0; i-- {
        oldName := rf.basePath
        if i > 1 {
            oldName = rf.basePath + "." + strconv.Itoa(i-1)
        }
        newName := rf.basePath + "." + strconv.Itoa(i)
        
        if _, err := os.Stat(oldName); err == nil {
            if err := os.Rename(oldName, newName); err != nil {
                return err
            }
        }
    }
    
    return rf.openFile()
}

func (rf *RotatingFile) Write(p []byte) (int, error) {
    rf.mu.Lock()
    defer rf.mu.Unlock()
    
    if rf.currentSize+int64(len(p)) > rf.maxSize {
        if err := rf.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := rf.file.Write(p)
    if err == nil {
        rf.currentSize += int64(n)
    }
    return n, err
}

func (rf *RotatingFile) Close() error {
    rf.mu.Lock()
    defer rf.mu.Unlock()
    
    if rf.file != nil {
        return rf.file.Close()
    }
    return nil
}

func main() {
    logFile, err := NewRotatingFile("app.log", 1024*1024) // 1MB max size
    if err != nil {
        fmt.Printf("Failed to create log file: %v\n", err)
        return
    }
    defer logFile.Close()
    
    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
        if _, err := logFile.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
    }
    
    fmt.Println("Log rotation test completed")
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

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath:  basePath,
        maxSize:   maxSize,
        fileIndex: 0,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    filename := fmt.Sprintf("%s.%d.log", l.basePath, l.fileIndex)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = info.Size()
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

        oldPath := fmt.Sprintf("%s.%d.log", l.basePath, l.fileIndex)
        compressedPath := fmt.Sprintf("%s.%d.log.gz", l.basePath, l.fileIndex)

        if err := compressFile(oldPath, compressedPath); err != nil {
            return err
        }

        if err := os.Remove(oldPath); err != nil {
            return err
        }

        l.fileIndex++
        l.currentSize = 0
    }

    return l.openCurrentFile()
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    target, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer target.Close()

    gz := gzip.NewWriter(target)
    defer gz.Close()

    _, err = io.Copy(gz, source)
    return err
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
    logger, err := NewRotatingLogger("app", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}