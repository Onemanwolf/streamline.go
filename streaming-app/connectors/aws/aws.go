package aws



import (
    "context"
    "github.com/segmentio/kafka-go"
    "streamline.go/config"
	"fmt"
)

type AWSConnector struct {
    writer *kafka.Writer
    config *config.AWSConfig
}

func NewAWSConnector(cfg *config.AWSConfig) (*AWSConnector, error) {
    return &AWSConnector{
        config: cfg,
    }, nil
}

func (a *AWSConnector) Connect() error {
    // For MSK, we use the Kafka protocol but with AWS-specific bootstrap servers
    // Note: In a real implementation, you'd need to fetch the bootstrap servers
    // from AWS MSK using the AWS SDK based on the stream name, region, etc.
    // Here we'll assume the brokers would be provided or retrieved separately

    // For this example, we'll simulate MSK connection using kafka-go
    a.writer = &kafka.Writer{
        Addr:     kafka.TCP("msk-broker-endpoint:9092"), // This would come from AWS MSK configuration
        Topic:    a.config.Stream,                       // Using stream name as topic
        Balancer: &kafka.LeastBytes{},
    }
    return nil
}

func (a *AWSConnector) Send(message string) error {
    if a.writer == nil {
        return fmt.Errorf("connector not initialized - call Connect() first")
    }
    return a.writer.WriteMessages(context.Background(),
        kafka.Message{
            Value: []byte(message),
        },
    )
}

func (a *AWSConnector) Close() error {
    if a.writer == nil {
        return nil
    }
    return a.writer.Close()
}