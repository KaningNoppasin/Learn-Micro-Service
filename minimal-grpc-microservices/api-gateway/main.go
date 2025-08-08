package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "api-gateway/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userClient pb.UserServiceClient

func main() {
	// Connect to user service
	// conn, err := grpc.Dial("user-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer conn.Close()

	userClient = pb.NewUserServiceClient(conn)

	r := gin.Default()

	// Routes
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)

	log.Println("API Gateway running on :8080")
	r.Run(":8080")
}

func createUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    response.Id,
		"name":  response.Name,
		"email": response.Email,
	})
}

func getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := userClient.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    response.Id,
		"name":  response.Name,
		"email": response.Email,
	})
}
