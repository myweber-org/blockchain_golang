package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

type AppConfig struct {
    ServerPort int
    DatabaseURL string
    DebugMode bool
    AllowedHosts []string
}

func LoadConfig() (*AppConfig, error) {
    cfg := &AppConfig{}
    
    portStr := os.Getenv("SERVER_PORT")
    if portStr == "" {
        portStr = "8080"
    }
    
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %v", err)
    }
    cfg.ServerPort = port
    
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }
    cfg.DatabaseURL = dbURL
    
    debugStr := os.Getenv("DEBUG_MODE")
    cfg.DebugMode = strings.ToLower(debugStr) == "true"
    
    hostsStr := os.Getenv("ALLOWED_HOSTS")
    if hostsStr != "" {
        cfg.AllowedHosts = strings.Split(hostsStr, ",")
    } else {
        cfg.AllowedHosts = []string{"localhost", "127.0.0.1"}
    }
    
    return cfg, nil
}

func (c *AppConfig) Validate() error {
    if c.ServerPort < 1 || c.ServerPort > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }
    
    if !strings.HasPrefix(c.DatabaseURL, "postgres://") && 
       !strings.HasPrefix(c.DatabaseURL, "mysql://") {
        return fmt.Errorf("database URL must start with postgres:// or mysql://")
    }
    
    return nil
}