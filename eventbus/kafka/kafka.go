package kafka

import (
	"context"
	"log/slog"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
)

// EventBus — Kafka реализация шины событий.
// Композиция Producer и Consumer.
type EventBus struct {
	producer *Producer
	consumer *Consumer
}

// New создаёт Kafka шину событий.
func New(options Options, logger *slog.Logger) *EventBus {
	return &EventBus{
		producer: NewProducer(options, logger),
		consumer: NewConsumer(options, logger),
	}
}

// Publish публикует доменное событие через Kafka producer.
func (b *EventBus) Publish(ctx context.Context, event *events.DomainEvent) error {
	return b.producer.Publish(ctx, event)
}

// PublishAsync публикует доменное событие через Kafka producer.
func (b *EventBus) PublishAsync(ctx context.Context, event *events.DomainEvent) error {
	return b.producer.Publish(ctx, event)
}

// Subscribe подписывает обработчик на события ограниченного контекста.
func (b *EventBus) Subscribe(boundedContext string, handler events.DomainEventHandler) {
	b.consumer.Subscribe(boundedContext, handler)
}

// Unsubscribe отписывает обработчик от событий ограниченного контекста.
func (b *EventBus) Unsubscribe(boundedContext string, handler events.DomainEventHandler) {
	b.consumer.Unsubscribe(boundedContext, handler)
}

// StartConsuming запускает background consumer goroutines.
func (b *EventBus) StartConsuming(ctx context.Context) {
	b.consumer.StartConsuming(ctx)
}

// Stop останавливает все consumer goroutines.
func (b *EventBus) Stop() {
	b.consumer.Stop()
}
