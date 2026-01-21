package config

import (
    "fmt"
    "os"
    "path/filepath"

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
    setFieldFromEnv(&config.Server.ReadTimeout, "READ_TIMEOUT")
    setFieldFromEnv(&config.Server.WriteTimeout, "WRITE_TIMEOUT")
    setFieldFromEnv(&config.Server.DebugMode, "DEBUG_MODE")
    
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
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".app", "config.yaml")
}package config

import (
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
    Logging struct {
        Level  string `yaml:"level"`
        Output string `yaml:"output"`
    } `yaml:"logging"`
}

func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

func (c *Config) Validate() error {
    if c.Server.Host == "" {
        c.Server.Host = "localhost"
    }
    if c.Server.Port == 0 {
        c.Server.Port = 8080
    }
    if c.Logging.Level == "" {
        c.Logging.Level = "info"
    }
    return nil
}