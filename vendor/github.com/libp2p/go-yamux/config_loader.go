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
        LogLevel: getEnv("LOG_LEVEL", "info"),
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
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    if value, err := strconv.Atoi(strValue); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    strValue := getEnv(key, "")
    if strValue == "" {
        return defaultValue
    }
    strValue = strings.ToLower(strValue)
    return strValue == "true" || strValue == "1" || strValue == "yes"
}

func validateConfig(config *Config) error {
    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return &ConfigError{Field: "DB_PORT", Message: "port must be between 1 and 65535"}
    }
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return &ConfigError{Field: "SERVER_PORT", Message: "port must be between 1 and 65535"}
    }
    if config.Server.ReadTimeout <= 0 {
        return &ConfigError{Field: "READ_TIMEOUT", Message: "timeout must be positive"}
    }
    if config.Server.WriteTimeout <= 0 {
        return &ConfigError{Field: "WRITE_TIMEOUT", Message: "timeout must be positive"}
    }
    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[strings.ToLower(config.LogLevel)] {
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
}