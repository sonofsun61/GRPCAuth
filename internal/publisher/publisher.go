package publisher

import (
	"authservice/internal/events"
	"authservice/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

var _ service.EventPublisher = (*Publisher)(nil)

type Publisher struct {
	writer  *kafka.Writer
	timeout time.Duration
}

func NewPublisher(brokers []string, topic string) *Publisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}
	return &Publisher{
		writer:  writer,
		timeout: 2 * time.Second,
	}
}

func (p *Publisher) PublishUserRegistered(ctx context.Context, eventInfo events.UserRegistered) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	eventInBytes, err := json.Marshal(eventInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(eventInfo.UserID.String()),
		Value: eventInBytes,
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish event to kafka: %w", err)
	}
	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
