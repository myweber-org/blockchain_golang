package main

import (
	"log"
	"os"
	"path/filepath"
)

const maxLogSize = 1024 * 1024 // 1MB

type RotatingLogger struct {
	file     *os.File
	filePath string
	baseName string
	counter  int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	baseName := base[:len(base)-len(ext)]

	rl := &RotatingLogger{
		filePath: path,
		baseName: baseName,
		counter:  0,
	}

	err := rl.openLogFile()
	return rl, err
}

func (rl *RotatingLogger) openLogFile() error {
	if rl.file != nil {
		rl.file.Close()
	}

	var err error
	rl.file, err = os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

func (rl *RotatingLogger) rotate() error {
	rl.counter++
	newPath := rl.filePath + "." + string(rune('0'+rl.counter))

	err := os.Rename(rl.filePath, newPath)
	if err != nil {
		return err
	}

	return rl.openLogFile()
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	info, err := rl.file.Stat()
	if err != nil {
		return 0, err
	}

	if info.Size()+int64(len(p)) > maxLogSize {
		err = rl.rotate()
		if err != nil {
			return 0, err
		}
	}

	return rl.file.Write(p)
}

func (rl *RotatingLogger) Close() error {
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)
	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry number %d", i)
	}
}