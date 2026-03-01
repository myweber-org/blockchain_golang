package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
    Name     string `yaml:"name" env:"DB_NAME"`
}

type ServerConfig struct {
    Port         int    `yaml:"port" env:"SERVER_PORT"`
    ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
    Debug        bool   `yaml:"debug" env:"SERVER_DEBUG"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    overrideFromEnv(&config)

    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    setFieldFromEnv(&config.Database.Host, "DB_HOST")
    setFieldFromEnv(&config.Database.Port, "DB_PORT")
    setFieldFromEnv(&config.Database.Username, "DB_USER")
    setFieldFromEnv(&config.Database.Password, "DB_PASS")
    setFieldFromEnv(&config.Database.Name, "DB_NAME")
    
    setFieldFromEnv(&config.Server.Port, "SERVER_PORT")
    setFieldFromEnv(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT")
    setFieldFromEnv(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT")
    setFieldFromEnv(&config.Server.Debug, "SERVER_DEBUG")
    
    setFieldFromEnv(&config.LogLevel, "LOG_LEVEL")
}

func setFieldFromEnv(field interface{}, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        switch v := field.(type) {
        case *string:
            *v = val
        case *int:
            fmt.Sscanf(val, "%d", v)
        case *bool:
            *v = val == "true" || val == "1"
        }
    }
}

func DefaultConfigPath() string {
    if path := os.Getenv("CONFIG_PATH"); path != "" {
        return path
    }
    
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./config.yaml"
    }
    
    return filepath.Join(homeDir, ".app", "config.yaml")
}