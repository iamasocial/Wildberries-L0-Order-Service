package kafka

import (
	"L0/internal/entities"
	"L0/internal/service"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	service service.Order
	reader  *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string, service service.Order) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		service: service,
		reader:  reader,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Println("Kafka consumer started")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		log.Printf("received message at topic/partion/offset %v/%v/%v: %s",
			m.Topic, m.Partition, m.Offset, string(m.Value))

		var order entities.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("failed to unmarshal order: %v", err)
			continue
		}

		if err := c.service.SaveOrder(ctx, &order); err != nil {
			log.Printf("failed to save order: %v", err)
			continue
		}

		log.Printf("order %s saved successfully", order.OrderUID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
