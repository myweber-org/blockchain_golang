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

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
    var cfg Config

    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("invalid config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&cfg)

    return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
    if val := os.Getenv("DB_HOST"); val != "" {
        cfg.Database.Host = val
    }
    if val := os.Getenv("DB_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Database.Port)
    }
    if val := os.Getenv("DB_USER"); val != "" {
        cfg.Database.Username = val
    }
    if val := os.Getenv("DB_PASS"); val != "" {
        cfg.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        cfg.Database.Name = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.Port)
    }
    if val := os.Getenv("READ_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.ReadTimeout)
    }
    if val := os.Getenv("WRITE_TIMEOUT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.WriteTimeout)
    }
    if val := os.Getenv("DEBUG_MODE"); val != "" {
        cfg.Server.DebugMode = val == "true" || val == "1"
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        cfg.LogLevel = val
    }
}package config

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/yaml.v2"
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

func LoadConfig(configPath string) (*Config, error) {
    config := &Config{}

    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    decoder := yaml.NewDecoder(file)
    if err := decoder.Decode(config); err != nil {
        return nil, fmt.Errorf("failed to parse config file: %w", err)
    }

    overrideFromEnv(config)

    return config, nil
}

func overrideFromEnv(c *Config) {
    if val := os.Getenv("SERVER_HOST"); val != "" {
        c.Server.Host = val
    }
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &c.Server.Port)
    }
    if val := os.Getenv("DB_HOST"); val != "" {
        c.Database.Host = val
    }
    if val := os.Getenv("DB_USERNAME"); val != "" {
        c.Database.Username = val
    }
    if val := os.Getenv("DB_PASSWORD"); val != "" {
        c.Database.Password = val
    }
    if val := os.Getenv("DB_NAME"); val != "" {
        c.Database.Name = val
    }
    if val := os.Getenv("LOG_LEVEL"); val != "" {
        c.LogLevel = strings.ToUpper(val)
    }
}