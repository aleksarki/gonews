package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type SaveServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type NewsAPIConfig struct {
	APIKey string `yaml:"api_key"`
}

type Config struct {
	Redis       RedisConfig       `yaml:"redis"`
	GRPC        GRPCConfig        `yaml:"grpc"`
	SaveService SaveServiceConfig `yaml:"save_service"`
	NewsAPI     NewsAPIConfig     `yaml:"newsapi"`
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
