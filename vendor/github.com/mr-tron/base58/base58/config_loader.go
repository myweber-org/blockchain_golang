package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env:"SERVER_HOST"`
		Port int    `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		URL      string `yaml:"url" env:"DB_URL"`
		MaxConns int    `yaml:"max_conns" env:"DB_MAX_CONNS"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
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
	setFieldFromEnv(&cfg.Server.Host, "SERVER_HOST")
	setFieldFromEnv(&cfg.Server.Port, "SERVER_PORT")
	setFieldFromEnv(&cfg.Database.URL, "DB_URL")
	setFieldFromEnv(&cfg.Database.MaxConns, "DB_MAX_CONNS")
	setFieldFromEnv(&cfg.LogLevel, "LOG_LEVEL")
}

func setFieldFromEnv(field interface{}, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		switch v := field.(type) {
		case *string:
			*v = val
		case *int:
			if intVal, err := parseInt(val); err == nil {
				*v = intVal
			}
		}
	}
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}package config

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
		return nil, errors.New("invalid database port")
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
		return nil, errors.New("invalid server port")
	}

	readTimeout, err := strconv.Atoi(getEnv("READ_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid read timeout")
	}

	writeTimeout, err := strconv.Atoi(getEnv("WRITE_TIMEOUT", "30"))
	if err != nil {
		return nil, errors.New("invalid write timeout")
	}

	debugMode := strings.ToLower(getEnv("DEBUG_MODE", "false")) == "true"

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
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return errors.New("database host cannot be empty")
	}

	if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
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