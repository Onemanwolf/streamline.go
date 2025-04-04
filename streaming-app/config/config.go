package config

import (
    "encoding/json"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Kafka  *KafkaConfig  `json:"kafka,omitempty" yaml:"kafka,omitempty"`
    Azure  *AzureConfig  `json:"azure,omitempty" yaml:"azure,omitempty"`
    AWS    *AWSConfig    `json:"aws,omitempty" yaml:"aws,omitempty"`
}

type KafkaConfig struct {
    Brokers []string `json:"brokers" yaml:"brokers"`
    Topic   string   `json:"topic" yaml:"topic"`
}

type AzureConfig struct {
    ConnectionString string `json:"connectionString" yaml:"connectionString"`
    EventHubName     string `json:"eventHubName" yaml:"eventHubName"`
}

type AWSConfig struct {
    AccessKey string `json:"accessKey" yaml:"accessKey"`
    SecretKey string `json:"secretKey" yaml:"secretKey"`
    Region    string `json:"region" yaml:"region"`
    Stream    string `json:"stream" yaml:"stream"`
}

func ParseConfig(data []byte) (*Config, error) {
    var config Config

    // Try JSON first
    if err := json.Unmarshal(data, &config); err == nil {
        return &config, nil
    }

    // Then try YAML
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}