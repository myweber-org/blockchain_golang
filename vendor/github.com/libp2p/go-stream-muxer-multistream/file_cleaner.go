package main

import (
    "os"
    "path/filepath"
    "time"
)

func cleanOldFiles(dir string, maxAgeDays int) error {
    cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if info.ModTime().Before(cutoff) {
            return os.Remove(path)
        }
        return nil
    })
}

func main() {
    tempDir := "/tmp"
    err := cleanOldFiles(tempDir, 7)
    if err != nil {
        panic(err)
    }
}