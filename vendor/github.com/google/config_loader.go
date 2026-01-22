package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path cannot be empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	err = validateConfig(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(c *Config) error {
	if c.Server.Host == "" {
		return errors.New("server host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.Database.Driver == "" {
		return errors.New("database driver is required")
	}
	if c.Database.Host == "" {
		return errors.New("database host is required")
	}
	if c.Database.Name == "" {
		return errors.New("database name is required")
	}
	return nil
}package config

import (
    "io/ioutil"
    "log"

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

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

func ValidateConfig(config *Config) bool {
    if config.Database.Host == "" {
        log.Println("Database host is required")
        return false
    }
    if config.Database.Port <= 0 {
        log.Println("Database port must be positive")
        return false
    }
    if config.Server.Port <= 0 {
        log.Println("Server port must be positive")
        return false
    }
    return true
}