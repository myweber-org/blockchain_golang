package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config map[string]string

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := make(Config)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid syntax at line %d", lineNumber)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			envVar := strings.TrimPrefix(strings.TrimSuffix(value, "}"), "${")
			envValue, exists := os.LookupEnv(envVar)
			if !exists {
				return nil, fmt.Errorf("environment variable %s not set at line %d", envVar, lineNumber)
			}
			value = envValue
		}

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return config, nil
}

func (c Config) Get(key string) (string, error) {
	value, exists := c[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in configuration", key)
	}
	return value, nil
}