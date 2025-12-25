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
}