package shared

import "os"

type Config struct {
	RabbitMQURL string
	ServerPort  string
}

func LoadConfig() *Config {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		RabbitMQURL: rabbitURL,
		ServerPort:  port,
	}
}
