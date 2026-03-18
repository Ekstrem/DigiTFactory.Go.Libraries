package events

import "context"

// DomainEventHandler описывает обработчик доменных событий.
type DomainEventHandler interface {
	// HandleEvent обрабатывает доменное событие.
	HandleEvent(ctx context.Context, event *DomainEvent) error
}
