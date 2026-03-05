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
	} `yaml:"database"`
	Logging struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
		File  string `yaml:"file" env:"LOG_FILE"`
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
	overrideString(&cfg.Server.Host, "SERVER_HOST")
	overrideInt(&cfg.Server.Port, "SERVER_PORT")
	overrideString(&cfg.Database.Host, "DB_HOST")
	overrideInt(&cfg.Database.Port, "DB_PORT")
	overrideString(&cfg.Database.Name, "DB_NAME")
	overrideString(&cfg.Database.User, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASSWORD")
	overrideString(&cfg.Logging.Level, "LOG_LEVEL")
	overrideString(&cfg.Logging.File, "LOG_FILE")
}

func overrideString(field *string, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		*field = val
	}
}

func overrideInt(field *int, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		var intVal int
		if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
			*field = intVal
		}
	}
}package config

import (
    "fmt"
    "io/ioutil"
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
    Port int    `yaml:"port"`
    Env  string `yaml:"env"`
}

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(path string) (*AppConfig, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
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
    if c.Server.Port <= 0 {
        return fmt.Errorf("server port must be positive")
    }
    return nil
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
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" env:"SERVER_PORT"`
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Version  string         `yaml:"version"`
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
	overrideString := func(field *string, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val
		}
	}

	overrideInt := func(field *int, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			var intVal int
			if _, err := fmt.Sscanf(val, "%d", &intVal); err == nil {
				*field = intVal
			}
		}
	}

	overrideBool := func(field *bool, envVar string) {
		if val := os.Getenv(envVar); val != "" {
			*field = val == "true" || val == "1" || val == "yes"
		}
	}

	overrideString(&config.Database.Host, "DB_HOST")
	overrideInt(&config.Database.Port, "DB_PORT")
	overrideString(&config.Database.Username, "DB_USER")
	overrideString(&config.Database.Password, "DB_PASS")
	overrideString(&config.Database.Database, "DB_NAME")
	overrideString(&config.Database.SSLMode, "DB_SSL_MODE")

	overrideInt(&config.Server.Port, "SERVER_PORT")
	overrideInt(&config.Server.ReadTimeout, "SERVER_READ_TIMEOUT")
	overrideInt(&config.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT")
	overrideBool(&config.Server.DebugMode, "SERVER_DEBUG")
	overrideString(&config.Server.LogLevel, "LOG_LEVEL")

	return nil
}

func validateConfig(config *AppConfig) error {
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if config.Database.Database == "" {
		return errors.New("database name is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Server.ReadTimeout < 0 {
		return errors.New("server read timeout cannot be negative")
	}
	if config.Server.WriteTimeout < 0 {
		return errors.New("server write timeout cannot be negative")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[config.Server.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}