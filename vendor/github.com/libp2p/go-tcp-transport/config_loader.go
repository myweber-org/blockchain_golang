package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" default:"8080"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DebugMode  bool   `env:"DEBUG_MODE" default:"false"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)
		
		envTag := structField.Tag.Get("env")
		defaultTag := structField.Tag.Get("default")
		
		var value string
		if envTag != "" {
			value = os.Getenv(envTag)
		}
		
		if value == "" && defaultTag != "" {
			value = defaultTag
		}
		
		if value == "" {
			return nil, fmt.Errorf("missing value for field %s", structField.Name)
		}
		
		if err := setFieldValue(field, value); err != nil {
			return nil, fmt.Errorf("invalid value for field %s: %w", structField.Name, err)
		}
	}
	
	return cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return errors.New("unsupported field type")
	}
	return nil
}

func (c *Config) Validate() error {
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if strings.TrimSpace(c.DBHost) == "" {
		return errors.New("database host cannot be empty")
	}
	if c.DBPort < 1 || c.DBPort > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}
	return nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}