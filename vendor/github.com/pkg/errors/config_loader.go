package config

import (
    "fmt"
    "io/ioutil"
    "os"

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
    Port         int    `yaml:"port"`
    ReadTimeout  int    `yaml:"read_timeout"`
    WriteTimeout int    `yaml:"write_timeout"`
}

type AppConfig struct {
    Environment string         `yaml:"environment"`
    Database    DatabaseConfig `yaml:"database"`
    Server      ServerConfig   `yaml:"server"`
}

func LoadConfig(filePath string) (*AppConfig, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return nil, fmt.Errorf("config file not found: %s", filePath)
    }

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %v", err)
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %v", err)
    }

    return &config, nil
}

func validateConfig(config *AppConfig) error {
    if config.Environment == "" {
        return fmt.Errorf("environment must be specified")
    }

    if config.Database.Host == "" {
        return fmt.Errorf("database host must be specified")
    }

    if config.Database.Port <= 0 || config.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1 and 65535")
    }

    if config.Server.Port <= 0 || config.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    return nil
}