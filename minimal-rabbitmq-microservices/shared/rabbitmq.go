package shared

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{conn: conn, channel: channel}, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

// Publish sends a message to a queue
func (r *RabbitMQ) Publish(queueName string, message interface{}) error {
	// Declare queue (creates if doesn't exist)
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send message
	return r.channel.Publish("", queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent, // Message survives server restart
		ContentType:  "application/json",
		Body:         body,
	})
}

// Consume listens for messages on a queue
func (r *RabbitMQ) Consume(queueName string, handler func([]byte) error) error {
	// Declare queue
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Start consuming
	msgs, err := r.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// Process messages
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			err := handler(msg.Body)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack(false, true) // Retry message
			} else {
				msg.Ack(false) // Acknowledge success
			}
		}
	}()

	log.Printf("Listening for messages on queue: %s", queueName)
	<-forever
	return nil
}
