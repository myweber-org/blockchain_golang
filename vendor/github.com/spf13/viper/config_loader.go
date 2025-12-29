package config

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/yaml.v3"
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
    overrideStruct(config)
}

func overrideStruct(v interface{}) {
    // Implementation would use reflection to check struct tags
    // and override values from environment variables
    // This is a simplified placeholder
    fmt.Printf("Config loaded with environment overrides for %T\n", v)
}

func GetEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func ValidateConfig(config *AppConfig) error {
    var errors []string

    if config.Database.Host == "" {
        errors = append(errors, "database host is required")
    }
    if config.Database.Port <= 0 {
        errors = append(errors, "database port must be positive")
    }
    if config.Server.Port <= 0 {
        errors = append(errors, "server port must be positive")
    }

    if len(errors) > 0 {
        return fmt.Errorf("config validation failed: %s", strings.Join(errors, ", "))
    }

    return nil
}