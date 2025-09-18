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

type Service3 struct {
	conn     *rabbitmq.Conn
	consumer *rabbitmq.Consumer
}

func NewService3() (*Service3, error) {
	// Create connection
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}

	// Create consumer for service3 queue
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"service3_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("service3.process"),
		rabbitmq.WithConsumerOptionsExchangeName("service3_exchange"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Service3{
		conn:     conn,
		consumer: consumer,
	}, nil
}

func (s *Service3) Close() {
	if s.consumer != nil {
		s.consumer.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Service3) processMessage(d rabbitmq.Delivery) rabbitmq.Action {
	var message shared.Message
	if err := json.Unmarshal(d.Body, &message); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return rabbitmq.NackDiscard
	}

	log.Printf("Service3 processing final message: %s", message.ID)

	// Simulate final processing
	time.Sleep(500 * time.Millisecond)

	// Add final processing result to message
	if message.Data == nil {
		message.Data = make(map[string]interface{})
	}
	message.Data["service3_processed"] = true
	message.Data["service3_timestamp"] = time.Now()
	message.Data["final_result"] = "Processing completed successfully"

	// Log the final result
	log.Printf("FINAL RESULT for message %s: %+v", message.ID, message.Data)

	return rabbitmq.Ack
}

func (s *Service3) startConsumer() error {
	return s.consumer.Run(s.processMessage)
}

func main() {
	service3, err := NewService3()
	if err != nil {
		log.Fatal("Failed to create Service3:", err)
	}
	defer service3.Close()

	// Start consumer in a goroutine
	go func() {
		if err := service3.startConsumer(); err != nil {
			log.Fatal("Failed to start consumer:", err)
		}
	}()

	// Start HTTP server for health checks
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "service3"})
	})

	log.Println("Service3 starting on :8083")
	r.Run(":8083")
}
