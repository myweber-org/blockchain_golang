package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    Port        int
    DatabaseURL string
    Debug       bool
    MaxWorkers  int
}

func Load() (*Config, error) {
    cfg := &Config{
        Port:        8080,
        DatabaseURL: "postgres://localhost:5432/app",
        Debug:       false,
        MaxWorkers:  4,
    }

    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid port: %v", err)
        }
        if port < 1 || port > 65535 {
            return nil, fmt.Errorf("port out of range: %d", port)
        }
        cfg.Port = port
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if debugStr := os.Getenv("APP_DEBUG"); debugStr != "" {
        debug, err := strconv.ParseBool(strings.ToLower(debugStr))
        if err != nil {
            return nil, fmt.Errorf("invalid debug value: %v", err)
        }
        cfg.Debug = debug
    }

    if workersStr := os.Getenv("MAX_WORKERS"); workersStr != "" {
        workers, err := strconv.Atoi(workersStr)
        if err != nil {
            return nil, fmt.Errorf("invalid max workers: %v", err)
        }
        if workers < 1 {
            return nil, fmt.Errorf("max workers must be positive")
        }
        cfg.MaxWorkers = workers
    }

    return cfg, nil
}

func (c *Config) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }
    if c.DatabaseURL == "" {
        return fmt.Errorf("database URL is required")
    }
    if c.MaxWorkers < 1 {
        return fmt.Errorf("max workers must be positive")
    }
    return nil
}