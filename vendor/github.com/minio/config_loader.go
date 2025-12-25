package config

import (
	"errors"
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
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
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

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.loadFromEnv(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) loadFromEnv() error {
	envVars := map[string]*string{
		"SERVER_HOST":     &c.Server.Host,
		"DB_HOST":         &c.Database.Host,
		"DB_NAME":         &c.Database.Name,
		"DB_USER":         &c.Database.User,
		"DB_PASSWORD":     &c.Database.Password,
		"LOG_LEVEL":       &c.LogLevel,
	}

	for envName, fieldPtr := range envVars {
		if val, exists := os.LookupEnv(envName); exists && val != "" {
			*fieldPtr = val
		}
	}

	if c.Server.Port == 0 {
		if portStr, exists := os.LookupEnv("SERVER_PORT"); exists && portStr != "" {
			port, err := parseInt(portStr)
			if err != nil {
				return err
			}
			c.Server.Port = port
		}
	}

	if c.Database.Port == 0 {
		if portStr, exists := os.LookupEnv("DB_PORT"); exists && portStr != "" {
			port, err := parseInt(portStr)
			if err != nil {
				return err
			}
			c.Database.Port = port
		}
	}

	return nil
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, errors.New("invalid integer value")
	}
	return result, nil
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
    Port         int            `yaml:"port"`
    ReadTimeout  int            `yaml:"read_timeout"`
    WriteTimeout int            `yaml:"write_timeout"`
    Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*ServerConfig, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config ServerConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(c *ServerConfig) error {
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Port)
    }

    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }

    if c.Database.Port <= 0 || c.Database.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", c.Database.Port)
    }

    if c.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }

    return nil
}