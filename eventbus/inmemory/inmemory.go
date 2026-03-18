// Package inmemory предоставляет In-Memory реализацию EventBus.
// Синхронный dispatch всем подписанным хендлерам в текущем процессе.
package inmemory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
)

// EventBus — In-Memory реализация шины событий.
type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]events.DomainEventHandler
	logger   *slog.Logger
}

// New создаёт In-Memory шину событий.
func New(logger *slog.Logger) *EventBus {
	return &EventBus{
		handlers: make(map[string][]events.DomainEventHandler),
		logger:   logger,
	}
}

// Publish публикует доменное событие синхронно всем подписчикам.
func (b *EventBus) Publish(ctx context.Context, event *events.DomainEvent) error {
	b.mu.RLock()
	handlers := b.handlers[event.BoundedContextName]
	// Копируем слайс для безопасного итерирования без удержания лока
	snapshot := make([]events.DomainEventHandler, len(handlers))
	copy(snapshot, handlers)
	b.mu.RUnlock()

	for _, handler := range snapshot {
		if err := handler.HandleEvent(ctx, event); err != nil {
			b.logger.Error("error dispatching domain event",
				"handler", handlerName(handler),
				"boundedContext", event.BoundedContextName,
				"error", err)
		}
	}
	return nil
}

// PublishAsync публикует доменное событие (в In-Memory реализации идентично Publish).
func (b *EventBus) PublishAsync(ctx context.Context, event *events.DomainEvent) error {
	return b.Publish(ctx, event)
}

// Subscribe подписывает обработчик на события указанного ограниченного контекста.
func (b *EventBus) Subscribe(boundedContext string, handler events.DomainEventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[boundedContext] = append(b.handlers[boundedContext], handler)
	b.logger.Debug("subscribed handler",
		"boundedContext", boundedContext)
}

// Unsubscribe отписывает обработчик от событий указанного ограниченного контекста.
func (b *EventBus) Unsubscribe(boundedContext string, handler events.DomainEventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[boundedContext]
	for i, h := range handlers {
		if h == handler {
			b.handlers[boundedContext] = append(handlers[:i], handlers[i+1:]...)
			b.logger.Debug("unsubscribed handler",
				"boundedContext", boundedContext)
			return
		}
	}
}

func handlerName(h events.DomainEventHandler) string {
	return "DomainEventHandler"
}
