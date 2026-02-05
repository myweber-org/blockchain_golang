package config

import (
    "bufio"
    "os"
    "strings"
)

type Config map[string]string

func LoadConfig(filename string) (Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    config := make(Config)
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        config[key] = os.ExpandEnv(value)
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return config, nil
}