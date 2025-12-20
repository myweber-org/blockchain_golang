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
    ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
    WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
    DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
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
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)
    return &config, nil
}

func overrideFromEnv(config *AppConfig) {
    setFromEnv(&config.Database.Host, "DB_HOST")
    setFromEnvInt(&config.Database.Port, "DB_PORT")
    setFromEnv(&config.Database.Username, "DB_USER")
    setFromEnv(&config.Database.Password, "DB_PASS")
    setFromEnv(&config.Database.Name, "DB_NAME")
    
    setFromEnvInt(&config.Server.Port, "SERVER_PORT")
    setFromEnvInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    setFromEnvInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    setFromEnvBool(&config.Server.DebugMode, "DEBUG_MODE")
    
    setFromEnv(&config.LogLevel, "LOG_LEVEL")
}

func setFromEnv(field *string, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val
    }
}

func setFromEnvInt(field *int, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        var intVal int
        if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
            *field = intVal
        }
    }
}

func setFromEnvBool(field *bool, envVar string) {
    if val := os.Getenv(envVar); val != "" {
        *field = val == "true" || val == "1" || val == "yes"
    }
}

func DefaultConfigPath() string {
    paths := []string{
        "config.yaml",
        "config.yml",
        filepath.Join("config", "config.yaml"),
        filepath.Join("config", "config.yml"),
        filepath.Join("..", "config", "config.yaml"),
    }
    
    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            return path
        }
    }
    
    return "config.yaml"
}