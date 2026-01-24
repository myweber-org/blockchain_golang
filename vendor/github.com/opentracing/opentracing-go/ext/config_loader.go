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
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	} `yaml:"database"`
	Logging struct {
		Level    string `yaml:"level" env:"LOG_LEVEL"`
		FilePath string `yaml:"file_path" env:"LOG_FILE_PATH"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
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

func overrideFromEnv(config *Config) error {
	envOverrides := map[string]func(string) error{
		"SERVER_HOST":    func(v string) error { config.Server.Host = v; return nil },
		"SERVER_PORT":    func(v string) error { return setInt(&config.Server.Port, v) },
		"DB_HOST":        func(v string) error { config.Database.Host = v; return nil },
		"DB_PORT":        func(v string) error { return setInt(&config.Database.Port, v) },
		"DB_NAME":        func(v string) error { config.Database.Name = v; return nil },
		"DB_USER":        func(v string) error { config.Database.User = v; return nil },
		"DB_PASSWORD":    func(v string) error { config.Database.Password = v; return nil },
		"DB_SSL_MODE":    func(v string) error { config.Database.SSLMode = v; return nil },
		"LOG_LEVEL":      func(v string) error { config.Logging.Level = v; return nil },
		"LOG_FILE_PATH":  func(v string) error { config.Logging.FilePath = v; return nil },
	}

	for envVar, setter := range envOverrides {
		if val, exists := os.LookupEnv(envVar); exists && val != "" {
			if err := setter(val); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateConfig(config *Config) error {
	if config.Server.Host == "" {
		return errors.New("server host cannot be empty")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Database.Host == "" {
		return errors.New("database host cannot be empty")
	}
	if config.Database.Name == "" {
		return errors.New("database name cannot be empty")
	}
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	return nil
}

func setInt(target *int, value string) error {
	var tmp int
	_, err := fmt.Sscanf(value, "%d", &tmp)
	if err != nil {
		return err
	}
	*target = tmp
	return nil
}