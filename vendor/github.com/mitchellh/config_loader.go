package config

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port    string `yaml:"port" env:"SERVER_PORT"`
		Timeout int    `yaml:"timeout" env:"SERVER_TIMEOUT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     string `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	Logging struct {
		Level  string `yaml:"level" env:"LOG_LEVEL"`
		Output string `yaml:"output" env:"LOG_OUTPUT"`
	} `yaml:"logging"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	if configPath != "" {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	overrideWithEnvVars(config)

	return config, nil
}

func overrideWithEnvVars(config *Config) {
	overrideStruct(config, "")
}

func overrideStruct(s interface{}, prefix string) {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			newPrefix := prefix
			if tag := fieldType.Tag.Get("yaml"); tag != "" {
				newPrefix = strings.ToUpper(tag) + "_"
			}
			overrideStruct(field.Addr().Interface(), newPrefix)
			continue
		}

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envVar := prefix + envTag
		if val := os.Getenv(envVar); val != "" {
			switch field.Kind() {
			case reflect.String:
				field.SetString(val)
			case reflect.Int:
				if intVal, err := strconv.Atoi(val); err == nil {
					field.SetInt(int64(intVal))
				}
			}
		}
	}
}