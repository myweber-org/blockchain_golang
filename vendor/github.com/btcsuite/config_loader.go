package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    CacheEnabled bool
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnvWithDefault("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT value: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := getEnvWithDefault("DATABASE_URL", "postgres://localhost:5432/app")
    if !strings.HasPrefix(dbURL, "postgres://") {
        return nil, fmt.Errorf("DATABASE_URL must start with postgres://")
    }
    cfg.DatabaseURL = dbURL
    
    cacheEnabled := getEnvWithDefault("CACHE_ENABLED", "true")
    cfg.CacheEnabled = strings.ToLower(cacheEnabled) == "true"
    
    logLevel := getEnvWithDefault("LOG_LEVEL", "info")
    validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLevels[strings.ToLower(logLevel)] {
        return nil, fmt.Errorf("invalid LOG_LEVEL: %s", logLevel)
    }
    cfg.LogLevel = strings.ToLower(logLevel)
    
    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
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
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", path)
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %v", err)
    }

    if config.Port == 0 {
        config.Port = 8080
    }
    if config.ReadTimeout == 0 {
        config.ReadTimeout = 30
    }
    if config.WriteTimeout == 0 {
        config.WriteTimeout = 30
    }

    return &config, nil
}

func (c *ServerConfig) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Port)
    }
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port < 1 || c.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", c.Database.Port)
    }
    return nil
}