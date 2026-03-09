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
}package config

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
    Debug        bool   `yaml:"debug" env:"DEBUG"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    LogLevel string         `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    overrideFromEnv(&config)
    return &config, nil
}

func overrideFromEnv(config *Config) {
    setFromEnv(&config.Database.Host, "DB_HOST")
    setFromEnvInt(&config.Database.Port, "DB_PORT")
    setFromEnv(&config.Database.Username, "DB_USER")
    setFromEnv(&config.Database.Password, "DB_PASS")
    setFromEnv(&config.Database.Name, "DB_NAME")

    setFromEnvInt(&config.Server.Port, "SERVER_PORT")
    setFromEnvInt(&config.Server.ReadTimeout, "READ_TIMEOUT")
    setFromEnvInt(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    setFromEnvBool(&config.Server.Debug, "DEBUG")

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
        "./config.yaml",
        "./config/config.yaml",
        "/etc/app/config.yaml",
    }

    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            return path
        }
    }
    return ""
}package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int            `yaml:"port"`
    Debug        bool           `yaml:"debug"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(filePath string) (*ServerConfig, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", filePath)
    }

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %v", err)
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    return &config, nil
}

func (c *ServerConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port < 1 || c.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    return nil
}