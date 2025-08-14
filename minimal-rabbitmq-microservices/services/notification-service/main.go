package main

import (
	"encoding/json"
	"log"
	"net/http"

	"minimal-rabbitmq-microservices/shared"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	config := shared.LoadConfig()

	// Connect to RabbitMQ
	rabbitmq, err := shared.NewRabbitMQ(config.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitmq.Close()

	// Setup HTTP server for health checks
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Start HTTP server in background
	go func() {
		log.Printf("Notification service HTTP starting on port %s", "8082")
		r.Run(":8082")
		// log.Printf("Notification service HTTP starting on port %s", config.ServerPort)
		// r.Run(":" + config.ServerPort)
	}()

	// Start consuming messages (this blocks)
	err = rabbitmq.Consume("order_notifications", handleOrderNotification)
	if err != nil {
		log.Fatal("Failed to consume messages:", err)
	}
}

func handleOrderNotification(body []byte) error {
	var event shared.OrderEvent

	// Parse message
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}

	// Process notification (simulate sending email/SMS)
	log.Printf("📧 Sending notification to User %d:", event.UserID)
	log.Printf("   Order #%d for %s ($%.2f)", event.OrderID, event.Product, event.Amount)
	log.Printf("   Message: %s", event.Message)

	// In real application, you would:
	// - Send actual email/SMS
	// - Save notification to database
	// - Call external notification service API

	return nil
}
