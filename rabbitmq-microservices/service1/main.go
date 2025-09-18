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

type Service1 struct {
	conn      *rabbitmq.Conn
	consumer  *rabbitmq.Consumer
	publisher *rabbitmq.Publisher
}

func NewService1() (*Service1, error) {
	// Create connection
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}

	// Create consumer for service1 queue
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"service1_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("service1.process"),
		rabbitmq.WithConsumerOptionsExchangeName("service1_exchange"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Create publisher for service2 exchange
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("service2_exchange"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		consumer.Close()
		conn.Close()
		return nil, err
	}

	return &Service1{
		conn:      conn,
		consumer:  consumer,
		publisher: publisher,
	}, nil
}

func (s *Service1) Close() {
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

func (s *Service1) processMessage(d rabbitmq.Delivery) rabbitmq.Action {
	var message shared.Message
	if err := json.Unmarshal(d.Body, &message); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return rabbitmq.NackDiscard
	}

	log.Printf("Service1 processing message: %s", message.ID)

	// Simulate processing
	time.Sleep(2 * time.Second)

	// Add processing result to message
	if message.Data == nil {
		message.Data = make(map[string]interface{})
	}
	message.Data["service1_processed"] = true
	message.Data["service1_timestamp"] = time.Now()
	message.Source = "service1"
	message.Step = 2

	// Send to Service2
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return rabbitmq.NackRequeue
	}

	err = s.publisher.Publish(
		messageBytes,
		[]string{"service2.process"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("service2_exchange"),
	)

	if err != nil {
		log.Printf("Failed to publish to service2: %v", err)
		return rabbitmq.NackRequeue
	}

	log.Printf("Service1 completed processing message: %s", message.ID)
	return rabbitmq.Ack
}

func (s *Service1) startConsumer() error {
	return s.consumer.Run(s.processMessage)
}

func main() {
	service1, err := NewService1()
	if err != nil {
		log.Fatal("Failed to create Service1:", err)
	}
	defer service1.Close()

	// Start consumer in a goroutine
	go func() {
		if err := service1.startConsumer(); err != nil {
			log.Fatal("Failed to start consumer:", err)
		}
	}()

	// Start HTTP server for health checks
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "service1"})
	})

	log.Println("Service1 starting on :8081")
	r.Run(":8081")
}
