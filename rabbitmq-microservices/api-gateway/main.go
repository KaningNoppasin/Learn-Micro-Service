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

type APIGateway struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
}

func NewAPIGateway() (*APIGateway, error) {
	// Create connection
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}

	// Create publisher for service1 exchange
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("service1_exchange"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &APIGateway{
		conn:      conn,
		publisher: publisher,
	}, nil
}

func (gw *APIGateway) Close() {
	if gw.publisher != nil {
		gw.publisher.Close()
	}
	if gw.conn != nil {
		gw.conn.Close()
	}
}

func (gw *APIGateway) processRequest(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create message for Service1
	message := shared.Message{
		ID:        generateID(),
		Data:      requestData,
		Timestamp: time.Now(),
		Source:    "api-gateway",
		Step:      1,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal message"})
		return
	}

	// Send to service1 with routing key
	err = gw.publisher.Publish(
		messageBytes,
		[]string{"service1.process"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("service1_exchange"),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "Request accepted and sent for processing",
		"request_id": message.ID,
	})
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().UnixNano()%1000))
}

func main() {
	gateway, err := NewAPIGateway()
	if err != nil {
		log.Fatal("Failed to create API Gateway:", err)
	}
	defer gateway.Close()

	r := gin.Default()

	r.POST("/process", gateway.processRequest)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	log.Println("API Gateway starting on :8080")
	r.Run(":8080")
}
