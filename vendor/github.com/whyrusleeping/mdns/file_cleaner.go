package main

import (
    "os"
    "path/filepath"
    "time"
)

const retentionDays = 7

func main() {
    tempDir := os.TempDir()
    cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

    filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info.IsDir() {
            return nil
        }
        if info.ModTime().Before(cutoffTime) {
            os.Remove(path)
        }
        return nil
    })
}