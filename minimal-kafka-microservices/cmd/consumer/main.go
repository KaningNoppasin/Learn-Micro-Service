package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"minimal-kafka/pkg/kafka"
	"minimal-kafka/pkg/models"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "events"
	port        = ":8081"
)

func main() {
	// Message handler function
	messageHandler := func(message interface{}) error {
		data, _ := json.Marshal(message)
		var event models.Event
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}

		log.Printf("Received event: ID=%s, Type=%s, Data=%s",
			event.ID, event.Type, event.Data)

		// Add your business logic here
		// Example: Save to database, call external API, etc.

		return nil
	}

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer([]string{kafkaBroker}, messageHandler)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}
	defer consumer.Close()

	// Start consumer in goroutine
	go func() {
		log.Println("Starting Kafka consumer...")
		if err := consumer.Start(topic); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Setup HTTP server for health checks
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Consumer is healthy",
		})
	})

	// Start HTTP server in goroutine
	go func() {
		log.Printf("Consumer HTTP server starting on %s", port)
		log.Fatal(r.Run(port))
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
