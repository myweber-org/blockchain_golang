package config

import (
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    dbConfig := DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnvAsInt("DB_PORT", 5432),
        Username: getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", ""),
        Database: getEnv("DB_NAME", "appdb"),
        SSLMode:  getEnv("DB_SSL_MODE", "disable"),
    }

    serverConfig := ServerConfig{
        Port:         getEnvAsInt("SERVER_PORT", 8080),
        ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
        WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
        DebugMode:    getEnvAsBool("DEBUG_MODE", false),
    }

    config := &Config{
        Database: dbConfig,
        Server:   serverConfig,
        LogLevel: strings.ToUpper(getEnv("LOG_LEVEL", "INFO")),
    }

    if err := validateConfig(config); err != nil {
        return nil, err
    }

    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
        return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
    }

    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
    }

    validLogLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
    if !validLogLevels[cfg.LogLevel] {
        return &ConfigError{Field: "LOG_LEVEL", Message: "invalid log level"}
    }

    return nil
}

type ConfigError struct {
    Field   string
    Message string
}

func (e *ConfigError) Error() string {
    return "config error: " + e.Field + " - " + e.Message
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

func LoadConfig(filePath string) (*ServerConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if config.Port == 0 {
        config.Port = 8080
    }

    return &config, nil
}

func ValidateConfig(config *ServerConfig) error {
    if config.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if config.Database.Port < 1 || config.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", config.Database.Port)
    }
    return nil
}

func GetEnvConfig() (*ServerConfig, error) {
    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        configPath = "config.yaml"
    }
    return LoadConfig(configPath)
}package config

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
    LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
    Version  string         `yaml:"version"`
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

    overrideFromEnv(&config.Database)
    overrideFromEnv(&config.Server)

    return &config, nil
}

func overrideFromEnv(config interface{}) {
    // Environment variable override logic would be implemented here
    // This is a placeholder for the actual implementation
}

func ValidateConfigPath(path string) error {
    absPath, err := filepath.Abs(path)
    if err != nil {
        return fmt.Errorf("invalid config path: %w", err)
    }

    info, err := os.Stat(absPath)
    if os.IsNotExist(err) {
        return fmt.Errorf("config file does not exist: %s", absPath)
    }
    if info.IsDir() {
        return fmt.Errorf("config path is a directory: %s", absPath)
    }

    ext := filepath.Ext(absPath)
    if ext != ".yaml" && ext != ".yml" {
        return fmt.Errorf("config file must be YAML format: %s", absPath)
    }

    return nil
}

func DefaultConfig() *Config {
    return &Config{
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Username: "postgres",
            Password: "",
            Name:     "appdb",
        },
        Server: ServerConfig{
            Port:         8080,
            ReadTimeout:  30,
            WriteTimeout: 30,
            DebugMode:    false,
            LogLevel:     "info",
        },
        Version: "1.0.0",
    }
}package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
}

type ServerConfig struct {
    Port         int    `yaml:"port"`
    ReadTimeout  int    `yaml:"read_timeout"`
    WriteTimeout int    `yaml:"write_timeout"`
}

type AppConfig struct {
    Environment string         `yaml:"environment"`
    Debug       bool           `yaml:"debug"`
    Database    DatabaseConfig `yaml:"database"`
    Server      ServerConfig   `yaml:"server"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML config: %w", err)
    }

    if config.Environment == "" {
        config.Environment = "development"
    }

    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }

    return &config, nil
}

func (c *AppConfig) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.Database.Port <= 0 {
        return fmt.Errorf("database port must be positive")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("database name is required")
    }
    return nil
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        Database: DatabaseConfig{
            Host:     getEnvWithDefault("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            Username: getEnvWithDefault("DB_USER", "postgres"),
            Password: getEnvWithDefault("DB_PASS", ""),
            Database: getEnvWithDefault("DB_NAME", "appdb"),
        },
        Server: ServerConfig{
            Port:         getEnvAsInt("SERVER_PORT", 8080),
            ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
            WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
            DebugMode:    getEnvAsBool("DEBUG_MODE", false),
        },
        LogLevel: getEnvWithDefault("LOG_LEVEL", "info"),
    }

    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }

    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }

    valueStr = strings.ToLower(valueStr)
    return valueStr == "true" || valueStr == "1" || valueStr == "yes"
}

func validateConfig(cfg *Config) error {
    if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", cfg.Database.Port)
    }

    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
    }

    if cfg.Server.ReadTimeout <= 0 {
        return fmt.Errorf("read timeout must be positive")
    }

    if cfg.Server.WriteTimeout <= 0 {
        return fmt.Errorf("write timeout must be positive")
    }

    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
        return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
    }

    return nil
}