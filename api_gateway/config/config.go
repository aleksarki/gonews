package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type SaveServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type SearchServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type NotifyServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	SaveService   SaveServiceConfig   `yaml:"save_service"`
	SearchService SearchServiceConfig `yaml:"search_service"`
	NotifyService NotifyServiceConfig `yaml:"notify_service"`
	HTTP          HTTPConfig          `yaml:"http"`
}

func LoadConfig(filename string) (*Config, error) {
	if filename == "" {
		filename = "config/config.yaml"
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &config, nil
}
