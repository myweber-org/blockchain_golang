package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	ServerPort string `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	LogLevel   string `json:"log_level"`
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			LogLevel:   getEnv("LOG_LEVEL", "info"),
		}

		configFile := os.Getenv("CONFIG_FILE")
		if configFile != "" {
			loadFromFile(configFile, instance)
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadFromFile(filename string, cfg *Config) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_ = decoder.Decode(cfg)
}package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	DebugMode  bool
	DatabaseURL string
	AllowedHosts []string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	
	portStr := getEnv("SERVER_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	cfg.ServerPort = port
	
	debugStr := getEnv("DEBUG_MODE", "false")
	cfg.DebugMode = strings.ToLower(debugStr) == "true"
	
	cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://localhost:5432/appdb")
	
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
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASS"`
	Database string `yaml:"database" env:"DB_NAME"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"DEBUG_MODE"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Features []string       `yaml:"features"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := overrideFromEnv(&config); err != nil {
		return nil, err
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func overrideFromEnv(config *AppConfig) error {
	envVars := map[string]string{
		"DB_HOST":       &config.Database.Host,
		"DB_PORT":       stringPtrFromInt(&config.Database.Port),
		"DB_USER":       &config.Database.Username,
		"DB_PASS":       &config.Database.Password,
		"DB_NAME":       &config.Database.Database,
		"SERVER_PORT":   stringPtrFromInt(&config.Server.Port),
		"READ_TIMEOUT":  stringPtrFromInt(&config.Server.ReadTimeout),
		"WRITE_TIMEOUT": stringPtrFromInt(&config.Server.WriteTimeout),
		"DEBUG_MODE":    stringPtrFromBool(&config.Server.DebugMode),
		"LOG_LEVEL":     &config.Server.LogLevel,
	}

	for envKey, fieldPtr := range envVars {
		if val, exists := os.LookupEnv(envKey); exists && val != "" {
			*fieldPtr = val
		}
	}

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("write timeout cannot be negative")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}

func stringPtrFromInt(i *int) *string {
	s := ""
	return &s
}

func stringPtrFromBool(b *bool) *string {
	s := ""
	return &s
}