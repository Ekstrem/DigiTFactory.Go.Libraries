package postgres

import (
	"context"
	"log/slog"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EventBus — PostgreSQL реализация шины событий (outbox pattern).
// Композиция Producer и Consumer.
type EventBus struct {
	producer *Producer
	consumer *Consumer
}

// New создаёт PostgreSQL шину событий.
func New(pool *pgxpool.Pool, options Options, logger *slog.Logger) *EventBus {
	return &EventBus{
		producer: NewProducer(pool, options, logger),
		consumer: NewConsumer(pool, options, logger),
	}
}

// Publish публикует доменное событие в outbox таблицу.
func (b *EventBus) Publish(ctx context.Context, event *events.DomainEvent) error {
	return b.producer.Publish(ctx, event)
}

// PublishAsync публикует доменное событие в outbox таблицу.
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

// StartConsuming запускает polling outbox таблицы.
func (b *EventBus) StartConsuming(ctx context.Context) {
	b.consumer.StartConsuming(ctx)
}

// Stop останавливает polling.
func (b *EventBus) Stop() {
	b.consumer.Stop()
}
