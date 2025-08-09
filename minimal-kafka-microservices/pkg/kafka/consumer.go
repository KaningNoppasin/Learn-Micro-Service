package kafka

import (
	"encoding/json"
	"log"

	// "github.com/Shopify/sarama"
	"github.com/IBM/sarama"
)

type MessageHandler func(message interface{}) error

type Consumer struct {
	consumer sarama.Consumer
	handler  MessageHandler
}

func NewConsumer(brokers []string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		handler:  handler,
	}, nil
}

func (c *Consumer) Start(topic string) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			var data interface{}
			if err := json.Unmarshal(message.Value, &data); err != nil {
				log.Printf("Failed to unmarshal: %v", err)
				continue
			}

			if err := c.handler(data); err != nil {
				log.Printf("Handler error: %v", err)
			}

		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
