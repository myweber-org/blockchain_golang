
package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type RotatingLogger struct {
	filePath    string
	maxSize     int64
	maxFiles    int
	currentSize int64
	file        *os.File
	logger      *log.Logger
}

func NewRotatingLogger(filePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: filePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
	rl.logger = log.New(file, "", log.LstdFlags)

	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()

	for i := rl.maxFiles - 1; i >= 0; i-- {
		oldPath := rl.getArchivePath(i)
		newPath := rl.getArchivePath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i+1 >= rl.maxFiles {
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	if err := os.Rename(rl.filePath, rl.getArchivePath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) getArchivePath(index int) string {
	if index == 0 {
		return rl.filePath + ".1"
	}
	return rl.filePath + "." + strconv.Itoa(index+1)
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	files, err := filepath.Glob(rl.filePath + ".*")
	if err != nil {
		return err
	}

	for _, file := range files {
		parts := strings.Split(file, ".")
		if len(parts) < 2 {
			continue
		}

		index, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			continue
		}

		if index > rl.maxFiles {
			os.Remove(file)
		}
	}

	return nil
}

func (rl *RotatingLogger) Println(v ...interface{}) {
	rl.logger.Println(v...)
}

func (rl *RotatingLogger) Close() error {
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		logger.Println("Log entry", i, "at", time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}

	logger.cleanupOldFiles()
}