package config

import (
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort    int
    DatabaseURL   string
    LogLevel      string
    CacheEnabled  bool
    MaxConnections int
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        ServerPort:    8080,
        DatabaseURL:   "localhost:5432",
        LogLevel:      "info",
        CacheEnabled:  true,
        MaxConnections: 100,
    }

    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
            cfg.ServerPort = port
        }
    }

    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        cfg.DatabaseURL = dbURL
    }

    if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
        validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
        if validLevels[strings.ToLower(logLevel)] {
            cfg.LogLevel = strings.ToLower(logLevel)
        }
    }

    if cacheFlag := os.Getenv("CACHE_ENABLED"); cacheFlag != "" {
        cfg.CacheEnabled = strings.ToLower(cacheFlag) == "true"
    }

    if maxConnStr := os.Getenv("MAX_CONNECTIONS"); maxConnStr != "" {
        if maxConn, err := strconv.Atoi(maxConnStr); err == nil && maxConn > 0 {
            cfg.MaxConnections = maxConn
        }
    }

    return cfg, nil
}