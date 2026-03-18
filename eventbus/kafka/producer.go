package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
	kafkago "github.com/segmentio/kafka-go"
)

// Producer — Kafka producer для доменных событий.
type Producer struct {
	options Options
	logger  *slog.Logger
}

// NewProducer создаёт Kafka producer.
func NewProducer(options Options, logger *slog.Logger) *Producer {
	return &Producer{
		options: options,
		logger:  logger,
	}
}

// Publish публикует доменное событие в Kafka.
func (p *Producer) Publish(ctx context.Context, event *events.DomainEvent) error {
	topic := p.topicName(event.BoundedContextName)
	data, err := Serialize(event)
	if err != nil {
		return fmt.Errorf("serialize event: %w", err)
	}

	writer := &kafkago.Writer{
		Addr:  kafkago.TCP(p.options.BootstrapServers),
		Topic: topic,
	}
	defer writer.Close()

	err = writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(event.CorrelationToken().String()),
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("publish to kafka topic %s: %w", topic, err)
	}

	p.logger.Debug("published event to kafka",
		"topic", topic,
		"aggregateId", event.AggregateID)
	return nil
}

func (p *Producer) topicName(boundedContext string) string {
	return fmt.Sprintf("%s.%s", p.options.TopicPrefix, strings.ToLower(boundedContext))
}
