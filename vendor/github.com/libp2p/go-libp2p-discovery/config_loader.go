package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    DatabaseURL string
    MaxConnections int
    DebugMode bool
    AllowedOrigins []string
}

func LoadConfig(filePath string) (*Config, error) {
    cfg := &Config{
        DatabaseURL: "localhost:5432",
        MaxConnections: 10,
        DebugMode: false,
        AllowedOrigins: []string{"http://localhost:3000"},
    }

    file, err := os.Open(filePath)
    if err != nil {
        return cfg, nil
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        switch key {
        case "DATABASE_URL":
            cfg.DatabaseURL = os.ExpandEnv(value)
        case "MAX_CONNECTIONS":
            if val, err := strconv.Atoi(value); err == nil {
                cfg.MaxConnections = val
            }
        case "DEBUG_MODE":
            cfg.DebugMode = strings.ToLower(value) == "true"
        case "ALLOWED_ORIGINS":
            cfg.AllowedOrigins = strings.Split(value, ",")
        }
    }

    return cfg, scanner.Err()
}