package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, group string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: group,
		}),
	}
}

func (c *Consumer) Read(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
