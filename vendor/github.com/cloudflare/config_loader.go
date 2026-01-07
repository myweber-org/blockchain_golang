
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" default:"8080"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     string `env:"DB_PORT" default:"5432"`
	DebugMode  bool   `env:"DEBUG_MODE" default:"false"`
	LogLevel   string `env:"LOG_LEVEL" default:"info"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	t := reflect.TypeOf(cfg).Elem()
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		defaultVal := field.Tag.Get("default")

		envValue := os.Getenv(envTag)
		if envValue == "" {
			envValue = defaultVal
		}

		fieldValue := v.Field(i)
		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Bool:
			boolVal := strings.ToLower(envValue) == "true"
			fieldValue.SetBool(boolVal)
		default:
			return nil, fmt.Errorf("unsupported field type: %s", field.Type.Kind())
		}
	}

	return cfg, nil
}

func (c *Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}