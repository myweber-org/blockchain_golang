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
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(path string) (*AppConfig, error) {
	config := &AppConfig{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func DefaultConfig() *AppConfig {
	config := &AppConfig{}
	config.Server.Host = "localhost"
	config.Server.Port = 8080
	config.Database.Host = "localhost"
	config.Database.Name = "appdb"
	config.Logging.Level = "info"
	config.Logging.Output = "stdout"
	return config
}

func ValidateConfig(config *AppConfig) bool {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		log.Printf("Invalid port number: %d", config.Server.Port)
		return false
	}
	if config.Database.Name == "" {
		log.Print("Database name cannot be empty")
		return false
	}
	return true
}package config

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
		PoolSize int    `yaml:"pool_size" env:"DB_POOL_SIZE"`
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
	overrideString(&cfg.Server.Host, "SERVER_HOST")
	overrideInt(&cfg.Server.Port, "SERVER_PORT")
	overrideString(&cfg.Database.URL, "DB_URL")
	overrideInt(&cfg.Database.PoolSize, "DB_POOL_SIZE")
	overrideString(&cfg.LogLevel, "LOG_LEVEL")
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
}