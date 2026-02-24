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
}package config

import (
    "fmt"
    "io"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    LogLevel string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &cfg, nil
}

func validateConfig(cfg *Config) error {
    if cfg.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    if cfg.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }

    return nil
}