package config

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
	ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	DebugMode    bool   `yaml:"debug_mode" env:"SERVER_DEBUG"`
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Features struct {
		Caching    bool `yaml:"caching" env:"FEATURE_CACHING"`
		Monitoring bool `yaml:"monitoring" env:"FEATURE_MONITORING"`
	} `yaml:"features"`
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

	if err := overrideFromEnv(&cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func overrideFromEnv(cfg *Config) error {
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

	overrideString(&cfg.Database.Host, "DB_HOST")
	overrideInt(&cfg.Database.Port, "DB_PORT")
	overrideString(&cfg.Database.Username, "DB_USER")
	overrideString(&cfg.Database.Password, "DB_PASS")
	overrideString(&cfg.Database.Database, "DB_NAME")

	overrideInt(&cfg.Server.Port, "SERVER_PORT")
	overrideInt(&cfg.Server.ReadTimeout, "SERVER_READ_TIMEOUT")
	overrideInt(&cfg.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT")
	overrideBool(&cfg.Server.DebugMode, "SERVER_DEBUG")
	overrideString(&cfg.Server.LogLevel, "LOG_LEVEL")

	overrideBool(&cfg.Features.Caching, "FEATURE_CACHING")
	overrideBool(&cfg.Features.Monitoring, "FEATURE_MONITORING")

	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return errors.New("database host is required")
	}
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[cfg.Server.LogLevel] {
		return errors.New("invalid log level")
	}
	return nil
}package config

import (
    "fmt"
    "io/ioutil"
    "os"

    "gopkg.in/yaml.v2"
)

type AppConfig struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    Database struct {
        Name     string `yaml:"name"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
    } `yaml:"database"`
    Logging struct {
        Level  string `yaml:"level"`
        Output string `yaml:"output"`
    } `yaml:"logging"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    configFile, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config AppConfig
    err = yaml.Unmarshal(configFile, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }

    return &config, nil
}

func ValidateConfig(config *AppConfig) error {
    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", config.Server.Port)
    }

    if config.Database.Name == "" {
        return fmt.Errorf("database name cannot be empty")
    }

    if config.Logging.Level != "debug" && config.Logging.Level != "info" && config.Logging.Level != "warn" && config.Logging.Level != "error" {
        return fmt.Errorf("invalid logging level: %s", config.Logging.Level)
    }

    return nil
}

func GetDefaultConfig() *AppConfig {
    var config AppConfig
    config.Server.Port = 8080
    config.Server.Host = "localhost"
    config.Database.Name = "appdb"
    config.Database.User = "admin"
    config.Database.Password = ""
    config.Database.Host = "localhost"
    config.Database.Port = 5432
    config.Logging.Level = "info"
    config.Logging.Output = "stdout"
    return &config
}

func SaveConfig(config *AppConfig, filePath string) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    err = ioutil.WriteFile(filePath, data, 0644)
    if err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }

    return nil
}

func LoadOrCreateConfig(filePath string) (*AppConfig, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        defaultConfig := GetDefaultConfig()
        err := SaveConfig(defaultConfig, filePath)
        if err != nil {
            return nil, fmt.Errorf("failed to create default config: %w", err)
        }
        return defaultConfig, nil
    }

    return LoadConfig(filePath)
}