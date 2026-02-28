package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    SSLMode  string
}

type ServerConfig struct {
    Port         int
    ReadTimeout  int
    WriteTimeout int
    DebugMode    bool
}

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    LogLevel string
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}

    dbHost := getEnvWithDefault("DB_HOST", "localhost")
    dbPort := getEnvAsInt("DB_PORT", 5432)
    dbUser := getEnvWithDefault("DB_USER", "postgres")
    dbPass := getEnvWithDefault("DB_PASS", "")
    dbName := getEnvWithDefault("DB_NAME", "appdb")
    dbSSL := getEnvWithDefault("DB_SSL_MODE", "disable")

    if dbPort < 1 || dbPort > 65535 {
        return nil, fmt.Errorf("invalid database port: %d", dbPort)
    }

    cfg.Database = DatabaseConfig{
        Host:     dbHost,
        Port:     dbPort,
        Username: dbUser,
        Password: dbPass,
        Database: dbName,
        SSLMode:  dbSSL,
    }

    serverPort := getEnvAsInt("SERVER_PORT", 8080)
    readTimeout := getEnvAsInt("READ_TIMEOUT", 30)
    writeTimeout := getEnvAsInt("WRITE_TIMEOUT", 30)
    debugMode := getEnvAsBool("DEBUG_MODE", false)

    if serverPort < 1 || serverPort > 65535 {
        return nil, fmt.Errorf("invalid server port: %d", serverPort)
    }

    cfg.Server = ServerConfig{
        Port:         serverPort,
        ReadTimeout:  readTimeout,
        WriteTimeout: writeTimeout,
        DebugMode:    debugMode,
    }

    logLevel := strings.ToLower(getEnvWithDefault("LOG_LEVEL", "info"))
    validLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLevels[logLevel] {
        return nil, fmt.Errorf("invalid log level: %s", logLevel)
    }

    cfg.LogLevel = logLevel

    return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnvWithDefault(key, "")
    if valueStr == "" {
        return defaultValue
    }

    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnvWithDefault(key, "")
    if valueStr == "" {
        return defaultValue
    }

    valueStr = strings.ToLower(valueStr)
    return valueStr == "true" || valueStr == "1" || valueStr == "yes"
}

func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("database host cannot be empty")
    }
    if c.Database.Username == "" {
        return fmt.Errorf("database username cannot be empty")
    }
    if c.Database.Database == "" {
        return fmt.Errorf("database name cannot be empty")
    }
    if c.Server.Port == 0 {
        return fmt.Errorf("server port cannot be zero")
    }
    return nil
}package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

func LoadConfig(filename string) (*AppConfig, error) {
	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	log.Printf("Configuration loaded successfully from %s", filename)
	return &config, nil
}

func ValidateConfig(config *AppConfig) bool {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return false
	}
	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return false
	}
	return true
}