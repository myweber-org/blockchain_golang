package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	overrideWithEnvVars(&config)
	return &config, nil
}

func getDefaultConfigPath() string {
	execPath, err := os.Executable()
	if err != nil {
		return "config.yaml"
	}
	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "config.yaml")
}

func overrideWithEnvVars(config *Config) {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = parseInt(port)
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = parseInt(port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.Username = user
	}
	if pass := os.Getenv("DB_PASS"); pass != "" {
		config.Database.Password = pass
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = strings.ToUpper(level)
	}
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		config.Logging.Output = output
	}
}

func parseInt(s string) int {
	var result int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int(ch-'0')
		} else {
			break
		}
	}
	return result
}

func ValidateConfig(config *Config) error {
	if config.Server.Host == "" {
		return errors.New("server host is required")
	}
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if config.Database.Host == "" {
		return errors.New("database host is required")
	}
	if config.Database.Name == "" {
		return errors.New("database name is required")
	}
	return nil
}package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `json:"server_port"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DebugMode  bool   `json:"debug_mode"`
	APIKeys    []string `json:"api_keys"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		ServerPort: 8080,
		DBHost:     "localhost",
		DBPort:     5432,
		DebugMode:  false,
		APIKeys:    []string{},
	}

	if configPath != "" {
		fileData, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(fileData, config); err != nil {
			return nil, err
		}
	}

	config.overrideFromEnv()

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) overrideFromEnv() {
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.ServerPort = port
		}
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		c.DBHost = host
	}

	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.DBPort = port
		}
	}

	if debugStr := os.Getenv("DEBUG_MODE"); debugStr != "" {
		c.DebugMode = strings.ToLower(debugStr) == "true"
	}

	if apiKeysStr := os.Getenv("API_KEYS"); apiKeysStr != "" {
		c.APIKeys = strings.Split(apiKeysStr, ",")
	}
}

func (c *Config) validate() error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return errors.New("invalid server port")
	}

	if c.DBPort < 1 || c.DBPort > 65535 {
		return errors.New("invalid database port")
	}

	if c.DBHost == "" {
		return errors.New("database host cannot be empty")
	}

	return nil
}