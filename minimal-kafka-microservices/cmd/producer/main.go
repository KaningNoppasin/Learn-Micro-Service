package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"minimal-kafka/pkg/kafka"
	"minimal-kafka/pkg/models"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "events"
	port        = ":8080"
)

func main() {
	// Initialize Kafka producer
	producer, err := kafka.NewProducer([]string{kafkaBroker})
	if err != nil {
		log.Fatal("Failed to create producer:", err)
	}
	defer producer.Close()

	// Setup Gin router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Producer is healthy",
		})
	})

	// Send event endpoint
	r.POST("/events", func(c *gin.Context) {
		var req struct {
			Type string `json:"type" binding:"required"`
			Data string `json:"data"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		event := models.Event{
			ID:      uuid.New().String(),
			Type:    req.Type,
			Data:    req.Data,
			Created: time.Now(),
		}

		if err := producer.Send(topic, event); err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "Failed to send event",
			})
			return
		}

		c.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Event sent successfully",
			Data:    event,
		})
	})

	log.Printf("Producer starting on %s", port)
	log.Fatal(r.Run(port))
}
