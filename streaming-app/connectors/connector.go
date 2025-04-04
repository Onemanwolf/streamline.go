package connectors

import (
	"streamline.go/config"
	"streamline.go/connectors/azure"
	"streamline.go/connectors/aws"
	"streamline.go/connectors/kafka"
	"fmt"
	





)

type Connector interface {
    Connect() error
    Send(message string) error
    Close() error
}

func NewConnector(provider string, configData []byte) (Connector, error) {
    cfg, err := config.ParseConfig(configData)
    if err != nil {
        return nil, err
    }

    switch provider {
    case "kafka":
        return kafka.NewKafkaConnector(cfg.Kafka)
    case "azure":
        return azure.NewAzureConnector(cfg.Azure)
    case "aws":
        return aws.NewAWSConnector(cfg.AWS)
    default:
        return nil, fmt.Errorf("unsupported provider: %s", provider)
    }
}