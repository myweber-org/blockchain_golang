
package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
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

	overrideWithEnv(&cfg)
	return &cfg, nil
}

func overrideWithEnv(cfg *Config) {
	setField(&cfg.Server.Host, "SERVER_HOST")
	setField(&cfg.Server.Port, "SERVER_PORT")
	setField(&cfg.Database.URL, "DB_URL")
	setField(&cfg.Database.PoolSize, "DB_POOL_SIZE")
	setField(&cfg.LogLevel, "LOG_LEVEL")
}

func setField(field interface{}, envVar string) {
	if val := os.Getenv(envVar); val != "" {
		switch v := field.(type) {
		case *string:
			*v = val
		case *int:
			if i, err := strconv.Atoi(val); err == nil {
				*v = i
			}
		}
	}
}

func (c *Config) GetDSN() string {
	return strings.Replace(c.Database.URL, "${DB_PASS}", os.Getenv("DB_PASSWORD"), 1)
}