package main

import (
	"log"
	"time"

	"minimal-rabbitmq-microservices/shared"

	"github.com/gin-gonic/gin"
)

var rabbitmq *shared.RabbitMQ

func main() {
	// Load config
	config := shared.LoadConfig()

	// Connect to RabbitMQ
	var err error
	rabbitmq, err = shared.NewRabbitMQ(config.RabbitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitmq.Close()

	// Setup Gin
	r := gin.Default()

	// Routes
	r.POST("/orders", createOrder)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	log.Printf("Order service starting on port %s", "8081")
	r.Run(":8081")
	// log.Printf("Order service starting on port %s", config.ServerPort)
	// r.Run(":" + config.ServerPort)
}

func createOrder(c *gin.Context) {
	var order shared.Order

	// Parse request
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Simulate saving order (in real app, save to database)
	order.ID = int(time.Now().Unix()) // Simple ID generation
	log.Printf("Created order: %+v", order)

	// Send event to notification service
	event := shared.OrderEvent{
		OrderID: order.ID,
		UserID:  order.UserID,
		Product: order.Product,
		Amount:  order.Amount,
		Message: "Your order has been created successfully!",
	}

	err := rabbitmq.Publish("order_notifications", event)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		c.JSON(500, gin.H{"error": "Order created but notification failed"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}
