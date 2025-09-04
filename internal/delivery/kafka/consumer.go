package kafka

import (
	"L0/internal/entities"
	"L0/internal/service"
	"L0/pkg/logger"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	service service.Order
	reader  *kafka.Reader
	logger  logger.Logger
}

func NewConsumer(brokers []string, topic, groupID string, service service.Order, logger logger.Logger) *Consumer {
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
		logger:  logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("Kafka consumer started")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.logger.Error("failed to read message from Kafka", "error", err.Error())
			return err
		}

		// log.Printf("received message at topic/partion/offset %v/%v/%v: %s",
		// 	m.Topic, m.Partition, m.Offset, string(m.Value))
		c.logger.Debug("received message from Kafka", "topic", m.Topic, "partition", m.Partition, "offset", m.Offset, "message", string(m.Value))

		var order entities.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			// log.Printf("failed to unmarshal order: %v", err)
			c.logger.Error("failed to unmarshal order", "error", err.Error())
			continue
		}

		if err := c.service.SaveOrder(ctx, &order); err != nil {
			// log.Printf("failed to save order: %v", err)
			c.logger.Error("failed to save order", "order_uid", order.OrderUID, "error", err.Error())
			continue
		}

		// log.Printf("order %s saved successfully", order.OrderUID)
		c.logger.Info("order saved successfully", "order_uid", order.OrderUID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
