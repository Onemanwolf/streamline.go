package kafka

import (
    "context"
    "log"
    "github.com/segmentio/kafka-go"
    "streamline.go/config"
	"fmt"
)

type KafkaConnector struct {
    writer *kafka.Writer
    config *config.KafkaConfig
}

func NewKafkaConnector(cfg *config.KafkaConfig) (*KafkaConnector, error) {
    return &KafkaConnector{
        config: cfg,
    }, nil
}

func (k *KafkaConnector) Connect() error {
    k.writer = &kafka.Writer{
        Addr:     kafka.TCP(k.config.Brokers...),
        Topic:    k.config.Topic,
        Balancer: &kafka.LeastBytes{},
        // Ensure messages are written synchronously for debugging
        BatchSize: 1,
        RequiredAcks: kafka.RequireAll, // Wait for all replicas to acknowledge
    }
    log.Printf("Connecting to Kafka brokers: %v, topic: %s", k.config.Brokers, k.config.Topic)
    return nil
}

func (k *KafkaConnector) Send(message string) error {
    if k.writer == nil {
        return fmt.Errorf("writer not initialized")
    }
    err := k.writer.WriteMessages(context.Background(),
        kafka.Message{
            Value: []byte(message),
        },
    )
    if err != nil {
        return fmt.Errorf("failed to write message: %v", err)
    }
    return nil
}

func (k *KafkaConnector) Close() error {
    if k.writer == nil {
        return nil
    }
    if err := k.writer.Close(); err != nil {
        return fmt.Errorf("failed to close writer: %v", err)
    }
    return nil
}