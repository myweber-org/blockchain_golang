package config

import (
	"errors"
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
	cfg := &Config{}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT value")
	}

	cfg.Database = DatabaseConfig{
		Host:     dbHost,
		Port:     dbPort,
		Username: getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASS", ""),
		Database: getEnv("DB_NAME", "appdb"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid SERVER_PORT value")
	}

	readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid READ_TIMEOUT value")
	}

	writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid WRITE_TIMEOUT value")
	}

	debugMode, err := strconv.ParseBool(getEnv("DEBUG_MODE", "false"))
	if err != nil {
		return nil, errors.New("invalid DEBUG_MODE value")
	}

	cfg.Server = ServerConfig{
		Port:         serverPort,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		DebugMode:    debugMode,
	}

	cfg.LogLevel = strings.ToUpper(getEnv("LOG_LEVEL", "INFO"))

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}

	validLogLevels := map[string]bool{
		"DEBUG": true,
		"INFO":  true,
		"WARN":  true,
		"ERROR": true,
		"FATAL": true,
	}

	if !validLogLevels[cfg.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}