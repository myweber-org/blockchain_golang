package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Format string `yaml:"format" env:"LOG_FORMAT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	overrideFromEnv(&cfg)
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	cfg.Server.Host = getEnvOrDefault("SERVER_HOST", cfg.Server.Host)
	cfg.Server.Port = getEnvIntOrDefault("SERVER_PORT", cfg.Server.Port)

	cfg.Database.Host = getEnvOrDefault("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvIntOrDefault("DB_PORT", cfg.Database.Port)
	cfg.Database.Name = getEnvOrDefault("DB_NAME", cfg.Database.Name)
	cfg.Database.User = getEnvOrDefault("DB_USER", cfg.Database.User)
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", cfg.Database.Password)
	cfg.Database.SSLMode = getEnvOrDefault("DB_SSL_MODE", cfg.Database.SSLMode)

	cfg.Logging.Level = getEnvOrDefault("LOG_LEVEL", cfg.Logging.Level)
	cfg.Logging.Format = getEnvOrDefault("LOG_FORMAT", cfg.Logging.Format)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    ServerPort int
    DBHost     string
    DBPort     int
    DebugMode  bool
    APIKeys    []string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    var err error
    cfg.ServerPort, err = getEnvInt("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    
    cfg.DBHost = getEnvString("DB_HOST", "localhost")
    
    cfg.DBPort, err = getEnvInt("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }
    
    cfg.DebugMode, err = getEnvBool("DEBUG_MODE", false)
    if err != nil {
        return nil, fmt.Errorf("invalid DEBUG_MODE: %w", err)
    }
    
    cfg.APIKeys = getEnvSlice("API_KEYS", []string{}, ",")
    
    if err := validateConfig(cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return cfg, nil
}

func getEnvString(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        intValue, err := strconv.Atoi(value)
        if err != nil {
            return 0, fmt.Errorf("cannot parse %s as integer: %w", key, err)
        }
        return intValue, nil
    }
    return defaultValue, nil
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
    if value := os.Getenv(key); value != "" {
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return false, fmt.Errorf("cannot parse %s as boolean: %w", key, err)
        }
        return boolValue, nil
    }
    return defaultValue, nil
}

func getEnvSlice(key string, defaultValue []string, separator string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, separator)
    }
    return defaultValue
}

func validateConfig(cfg *Config) error {
    if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    
    if cfg.DBPort < 1 || cfg.DBPort > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }
    
    if cfg.DBHost == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    
    return nil
}