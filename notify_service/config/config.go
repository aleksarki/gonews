package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type KafkaConfig struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	NotificationTopic string `yaml:"notification_topic"`
	ConsumerGroup     string `yaml:"consumer_group"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type SaveServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type SearchServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type SchedulerConfig struct {
	CheckIntervalMinutes int `yaml:"check_interval_minutes"`
}

type Config struct {
	Kafka         KafkaConfig         `yaml:"kafka"`
	GRPC          GRPCConfig          `yaml:"grpc"`
	SaveService   SaveServiceConfig   `yaml:"save_service"`
	SearchService SearchServiceConfig `yaml:"search_service"`
	Scheduler     SchedulerConfig     `yaml:"scheduler"`
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
