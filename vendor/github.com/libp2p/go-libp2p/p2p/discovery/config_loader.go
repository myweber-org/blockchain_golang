package config

import (
    "os"
    "strconv"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    EnableDebug bool
    MaxConnections int
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ServerPort:     8080,
        DatabaseURL:    "localhost:5432",
        EnableDebug:    false,
        MaxConnections: 10,
    }

    if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil {
            cfg.ServerPort = port
        }
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if debugStr := os.Getenv("ENABLE_DEBUG"); debugStr != "" {
        if debug, err := strconv.ParseBool(debugStr); err == nil {
            cfg.EnableDebug = debug
        }
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        if maxConn, err := strconv.Atoi(maxConnStr); err == nil {
            cfg.MaxConnections = maxConn
        }
    }

    return cfg, nil
}