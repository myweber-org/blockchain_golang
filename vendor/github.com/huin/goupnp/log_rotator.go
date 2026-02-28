
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
    file        *os.File
    currentSize int64
    maxSize     int64
    basePath    string
    counter     int
}

func NewRotatingLogger(basePath string, maxSizeMB int64) (*RotatingLogger, error) {
    if maxSizeMB <= 0 {
        return nil, fmt.Errorf("maxSize must be positive")
    }

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
        currentSize: info.Size(),
        maxSize:     maxSizeMB * 1024 * 1024,
        basePath:    basePath,
        counter:     0,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedName := fmt.Sprintf("%s.%s.%d.gz", rl.basePath, timestamp, rl.counter)
    rl.counter++

    if err := compressFile(rl.basePath, archivedName); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    return nil
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    gz := gzip.NewWriter(destination)
    defer gz.Close()

    _, err = io.Copy(gz, source)
    return err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type RotatingWriter struct {
	mu          sync.Mutex
	filename    string
	maxSize     int64
	currentSize int64
	file        *os.File
	rotation    int
}

func NewRotatingWriter(filename string, maxSize int64) (*RotatingWriter, error) {
	w := &RotatingWriter{
		filename: filename,
		maxSize:  maxSize,
	}

	if err := w.openFile(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *RotatingWriter) openFile() error {
	info, err := os.Stat(w.filename)
	if os.IsNotExist(err) {
		w.currentSize = 0
	} else if err != nil {
		return err
	} else {
		w.currentSize = info.Size()
	}

	file, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

func (w *RotatingWriter) rotate() error {
	w.file.Close()

	backupName := fmt.Sprintf("%s.%d", w.filename, w.rotation)
	if err := os.Rename(w.filename, backupName); err != nil {
		return err
	}

	w.rotation++
	w.currentSize = 0
	return w.openFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	if err != nil {
		return n, err
	}

	w.currentSize += int64(n)
	return n, nil
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}

func main() {
	writer, err := NewRotatingWriter("app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i)
		if _, err := writer.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}