package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(filename string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ValidateConfig(config *AppConfig) bool {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		log.Printf("Invalid server port: %d", config.Server.Port)
		return false
	}

	if config.Database.Host == "" {
		log.Print("Database host cannot be empty")
		return false
	}

	return true
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    Port        int
    DatabaseURL string
    LogLevel    string
    CacheTTL    int
}

func LoadConfig() (*AppConfig, error) {
    config := &AppConfig{
        Port:        getEnvAsInt("APP_PORT", 8080),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/app"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
        CacheTTL:    getEnvAsInt("CACHE_TTL", 300),
    }

    if err := validateConfig(config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
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

    value, err := strconv.Atoi(strValue)
    if err != nil {
        return defaultValue
    }
    return value
}

func validateConfig(config *AppConfig) error {
    if config.Port < 1 || config.Port > 65535 {
        return fmt.Errorf("invalid port number: %d", config.Port)
    }

    if !strings.HasPrefix(config.DatabaseURL, "postgres://") {
        return fmt.Errorf("invalid database URL format")
    }

    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLogLevels[config.LogLevel] {
        return fmt.Errorf("invalid log level: %s", config.LogLevel)
    }

    if config.CacheTTL < 0 {
        return fmt.Errorf("cache TTL cannot be negative")
    }

    return nil
}