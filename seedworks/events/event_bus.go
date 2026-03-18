package events

import "context"

// EventBusProducer описывает поставщика в шину событий.
type EventBusProducer interface {
	// Publish публикует доменное событие синхронно.
	Publish(ctx context.Context, event *DomainEvent) error

	// PublishAsync публикует доменное событие асинхронно.
	PublishAsync(ctx context.Context, event *DomainEvent) error
}

// EventBusConsumer описывает потребителя шины событий.
type EventBusConsumer interface {
	// Subscribe подписывается на доменные события указанного ограниченного контекста.
	Subscribe(boundedContext string, handler DomainEventHandler)

	// Unsubscribe отписывается от доменных событий указанного ограниченного контекста.
	Unsubscribe(boundedContext string, handler DomainEventHandler)
}

// EventBus описывает шину событий (объединяет Producer и Consumer).
type EventBus interface {
	EventBusProducer
	EventBusConsumer
}
