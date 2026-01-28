package config

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DatabaseURL string
    DebugMode bool
    AllowedHosts []string
}

func Load() (*Config, error) {
    cfg := &Config{}
    
    portStr := getEnv("SERVER_PORT", "8080")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, err
    }
    cfg.ServerPort = port
    
    cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/app")
    
    debugStr := getEnv("DEBUG_MODE", "false")
    cfg.DebugMode = strings.ToLower(debugStr) == "true"
    
    hostsStr := getEnv("ALLOWED_HOSTS", "localhost,127.0.0.1")
    cfg.AllowedHosts = strings.Split(hostsStr, ",")
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
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
    cfg := &Config{}

    dbHost := getEnvWithDefault("DB_HOST", "localhost")
    dbPort := getEnvIntWithDefault("DB_PORT", 5432)
    dbUser := getEnvWithDefault("DB_USER", "postgres")
    dbPass := getEnvWithDefault("DB_PASS", "")
    dbName := getEnvWithDefault("DB_NAME", "appdb")

    if dbPass == "" {
        return nil, fmt.Errorf("database password cannot be empty")
    }

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: dbUser,
        Password: dbPass,
        Database: dbName,
    }

    serverPort := getEnvIntWithDefault("SERVER_PORT", 8080)
    readTimeout := getEnvIntWithDefault("READ_TIMEOUT", 30)
    writeTimeout := getEnvIntWithDefault("WRITE_TIMEOUT", 30)
    debugMode := getEnvBoolWithDefault("DEBUG_MODE", false)

    if serverPort < 1 || serverPort > 65535 {
        return nil, fmt.Errorf("invalid server port: %d", serverPort)
    }

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
        DebugMode:    debugMode,
    }

    logLevel := strings.ToUpper(getEnvWithDefault("LOG_LEVEL", "INFO"))
    validLevels := map[string]bool{
        "DEBUG": true,
        "INFO":  true,
        "WARN":  true,
        "ERROR": true,
    }

    if !validLevels[logLevel] {
        return nil, fmt.Errorf("invalid log level: %s", logLevel)
    }

    cfg.LogLevel = logLevel

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvBoolWithDefault(key string, defaultValue bool) bool {
    if value, exists := os.LookupEnv(key); exists {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}