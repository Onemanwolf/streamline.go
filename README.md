# Streamline.go
## Streaming Connector App

## Overview

This application provides connectors for Kafka, Azure Event Hubs, and AWS MSK streaming services.

This app is a streaming message processing system that supports multiple connectors for different providers (Kafka, Azure Event Hubs, and AWS MSK). Here's an overview of its components:

1. **Connectors**:
   - The app defines a `Connector` interface with methods `Connect()`, `Send(message string)`, and `Close()`. Each provider (Kafka, Azure, AWS) implements this interface.
   - The `NewConnector` function dynamically creates a connector instance based on the provider name and configuration.

2. **Configuration**:
   - The `config` package parses configuration data in JSON or YAML format. It supports provider-specific configurations (Kafka, Azure, AWS).

3. **Provider Implementations**:
   - **Kafka**: Uses the `kafka-go` library to connect to Kafka brokers, send messages, and close the connection.
   - **Azure**: Uses the Azure Event Hubs SDK to connect to an Event Hub, send messages, and close the hub.
   - **AWS**: Simulates AWS MSK (Managed Streaming for Apache Kafka) using `kafka-go`. It connects to Kafka brokers and sends messages.

4. **HTTP API**:
   - The `main.go` file defines an HTTP API using the Gorilla Mux router.
   - The `/process` endpoint accepts a POST request with a JSON payload containing:
     - `provider`: The name of the provider (e.g., "kafka", "azure", "aws").
     - `configPath` or `config`: Configuration data for the provider.
     - `messages`: A list of messages to send.
   - The endpoint processes the messages asynchronously using a goroutine and a `sync.WaitGroup`.

5. **Message Processing**:
   - The app reads the configuration, initializes the appropriate connector, connects to the provider, sends the messages, and closes the connection.

6. **Server**:
   - The app starts an HTTP server on port `8081` to handle incoming requests.

This app is designed to be extensible, allowing additional providers to be added by implementing the `Connector` interface. It is useful for scenarios where messages need to be sent to different streaming platforms based on user input.


# Patterns

This app implements both the **Strategy Pattern** and the **Factory Pattern**:

### **Strategy Pattern**
The app uses the Strategy Pattern by defining a `Connector` interface with methods like `Connect()`, `Send(message string)`, and `Close()`. Each provider (Kafka, Azure, AWS) implements this interface with its own specific behavior. This allows the app to dynamically select and use a specific strategy (connector implementation) at runtime based on the provider.

### **Factory Pattern**
The app implements the Factory Pattern in the `NewConnector` function. This function acts as a factory that creates and returns an instance of the appropriate `Connector` implementation (Kafka, Azure, or AWS) based on the `provider` string. This encapsulates the object creation logic and abstracts it from the rest of the application.

By combining these patterns, the app achieves flexibility and extensibility, allowing new providers to be added with minimal changes to the existing code.

## HTTP Endpoint

POST /process

- Content-Type: application/json

- Body:

```  json
{
    "provider": "kafka|azure|aws",
    "configPath": "optional/path/to/config",
    "config": "optional JSON config object",
    "messages": ["message1", "message2"]
}

```

# Kafka Connector

## Configuration:

```json

{
    "kafka": {
        "brokers": ["localhost:29092"],
        "topic": "myTopic"
    }
}
```


```yaml

kafka:
    brokers:
        - localhost:29092
    topic: myTopic

```


```bash
curl -X POST http://localhost:8080/process \
-H "Content-Type: application/json" \
-d '{"provider":"kafka","config":{"kafka":{"brokers":["localhost:29092"],"topic":"myTopic"}},"messages":["hello world"]}'

```
# Azure Event Hubs Connector
## Configuration:

```yaml

```json

{
    "azure": {
        "connectionString": "<your-azure-event-hubs-connection-string>",
        "eventHubName": "<your-event-hub-name>"
    }
}

```

```bash
curl -X POST http://localhost:8081/process \
-H "Content-Type: application/json" \
-d '{"provider":"azure","config":{"azure":{"connectionString":"Endpoint=sb://myeventhubns.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=abc123xyz456...","eventHubName":"my-event-hub"}},"messages":["hello azure"]}'
```


# AWS MSK Connector
## Configuration:


```json

{
    "aws": {
        "accessKey": "<your-aws-access-key-id>",
        "secretKey": "<your-aws-secret-access-key>",
        "region": "<aws-region>",
        "stream": "<msk-cluster-name-or-topic>"
    }
}

```

```bash

curl -X POST http://localhost:8081/process \
-H "Content-Type: application/json" \
-d '{"provider":"aws","config":{"aws":{"accessKey":"AKIAIOSFODNN7EXAMPLE","secretKey":"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY","region":"us-west-2","stream":"my-msk-cluster"}},"messages":["hello aws"]}'


```