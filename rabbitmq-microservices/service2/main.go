package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "rabbitmq-microservices/shared"

    "github.com/gin-gonic/gin"
    "github.com/wagslane/go-rabbitmq"
)

type Service2 struct {
    conn      *rabbitmq.Conn
    consumer  *rabbitmq.Consumer
    publisher *rabbitmq.Publisher
}

func NewService2() (*Service2, error) {
    // Create connection
    conn, err := rabbitmq.NewConn(
        "amqp://guest:guest@localhost:5672/",
        rabbitmq.WithConnectionOptionsLogging,
    )
    if err != nil {
        return nil, err
    }

    // Create consumer for service2 queue
    consumer, err := rabbitmq.NewConsumer(
        conn,
        "service2_queue",
        rabbitmq.WithConsumerOptionsRoutingKey("service2.process"),
        rabbitmq.WithConsumerOptionsExchangeName("service2_exchange"),
        rabbitmq.WithConsumerOptionsExchangeDeclare,
    )
    if err != nil {
        conn.Close()
        return nil, err
    }

    // Create publisher for service3 exchange
    publisher, err := rabbitmq.NewPublisher(
        conn,
        rabbitmq.WithPublisherOptionsLogging,
        rabbitmq.WithPublisherOptionsExchangeName("service3_exchange"),
        rabbitmq.WithPublisherOptionsExchangeDeclare,
    )
    if err != nil {
        consumer.Close()
        conn.Close()
        return nil, err
    }

    return &Service2{
        conn:      conn,
        consumer:  consumer,
        publisher: publisher,
    }, nil
}

func (s *Service2) Close() {
    if s.consumer != nil {
        s.consumer.Close()
    }
    if s.publisher != nil {
        s.publisher.Close()
    }
    if s.conn != nil {
        s.conn.Close()
    }
}

func (s *Service2) processMessage(d rabbitmq.Delivery) rabbitmq.Action {
    var message shared.Message
    if err := json.Unmarshal(d.Body, &message); err != nil {
        log.Printf("Failed to unmarshal message: %v", err)
        return rabbitmq.NackDiscard
    }

    log.Printf("Service2 processing message: %s", message.ID)

    // Simulate processing
    time.Sleep(1 * time.Second)

    // Add processing result to message
    if message.Data == nil {
        message.Data = make(map[string]interface{})
    }
    message.Data["service2_processed"] = true
    message.Data["service2_timestamp"] = time.Now()
    message.Source = "service2"
    message.Step = 3

    // Send to Service3
    messageBytes, err := json.Marshal(message)
    if err != nil {
        log.Printf("Failed to marshal message: %v", err)
        return rabbitmq.NackRequeue
    }

    err = s.publisher.Publish(
        messageBytes,
        []string{"service3.process"},
        rabbitmq.WithPublishOptionsContentType("application/json"),
        rabbitmq.WithPublishOptionsExchange("service3_exchange"),
    )

    if err != nil {
        log.Printf("Failed to publish to service3: %v", err)
        return rabbitmq.NackRequeue
    }

    log.Printf("Service2 completed processing message: %s", message.ID)
    return rabbitmq.Ack
}

func (s *Service2) startConsumer() error {
    return s.consumer.Run(s.processMessage)
}

func main() {
    service2, err := NewService2()
    if err != nil {
        log.Fatal("Failed to create Service2:", err)
    }
    defer service2.Close()

    // Start consumer in a goroutine
    go func() {
        if err := service2.startConsumer(); err != nil {
            log.Fatal("Failed to start consumer:", err)
        }
    }()

    // Start HTTP server for health checks
    r := gin.Default()
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "service2"})
    })

    log.Println("Service2 starting on :8082")
    r.Run(":8082")
}
