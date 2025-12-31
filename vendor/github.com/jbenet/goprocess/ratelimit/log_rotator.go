package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
)

type RotatingWriter struct {
    currentSize int64
    file        *os.File
    basePath    string
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingWriter{
        currentSize: info.Size(),
        file:        file,
        basePath:    path,
    }, nil
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := w.file.Write(p)
    if err == nil {
        w.currentSize += int64(n)
    }
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if err := w.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", w.basePath, timestamp)

    if err := os.Rename(w.basePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(w.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    w.file = file
    w.currentSize = 0

    go w.cleanupOldBackups()
    return nil
}

func (w *RotatingWriter) cleanupOldBackups() {
    dir := filepath.Dir(w.basePath)
    baseName := filepath.Base(w.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if len(name) > len(baseName) && name[:len(baseName)] == baseName && name[len(baseName)] == '.' {
            backups = append(backups, name)
        }
    }

    if len(backups) > maxBackups {
        toRemove := backups[:len(backups)-maxBackups]
        for _, backup := range toRemove {
            os.Remove(filepath.Join(dir, backup))
        }
    }
}

func (w *RotatingWriter) Close() error {
    return w.file.Close()
}

func main() {
    writer, err := NewRotatingWriter("logs/application.log")
    if err != nil {
        panic(err)
    }
    defer writer.Close()

    for i := 0; i < 10000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        writer.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
}package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type RotatingWriter struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	rotationNum int
}

func NewRotatingWriter(filePath string, maxSize int64) (*RotatingWriter, error) {
	writer := &RotatingWriter{
		filePath: filePath,
		maxSize:  maxSize,
	}

	if err := writer.openCurrentFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *RotatingWriter) openCurrentFile() error {
	dir := filepath.Dir(w.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.currentFile = file
	w.currentSize = stat.Size()
	return nil
}

func (w *RotatingWriter) rotate() error {
	w.currentFile.Close()

	backupPath := fmt.Sprintf("%s.%d", w.filePath, w.rotationNum)
	if err := os.Rename(w.filePath, backupPath); err != nil {
		return err
	}

	w.rotationNum++
	return w.openCurrentFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.currentFile.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("logs/app.log", 1024*10)
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i)
		if _, err := writer.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}