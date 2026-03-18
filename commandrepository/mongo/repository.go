package mongo

import (
	"context"

	"github.com/google/uuid"
)

// EventStoreRepository описывает интерфейс MongoDB Event Store.
type EventStoreRepository interface {
	// GetByID получает все события агрегата по идентификатору.
	GetByID(ctx context.Context, id uuid.UUID) ([]DomainEventDocument, error)

	// GetByIDAndVersion получает событие по идентификатору и версии.
	GetByIDAndVersion(ctx context.Context, id uuid.UUID, version int64) (*DomainEventDocument, error)

	// GetByCorrelationToken получает событие по маркеру корреляции.
	GetByCorrelationToken(ctx context.Context, correlationToken uuid.UUID) (*DomainEventDocument, error)

	// SaveEvent сохраняет событие в Event Store.
	SaveEvent(ctx context.Context, doc *DomainEventDocument) error

	// SaveSnapshot сохраняет snapshot агрегата.
	SaveSnapshot(ctx context.Context, id uuid.UUID, version int64, aggregateJSON string) error

	// GetEventCount возвращает количество событий для агрегата.
	GetEventCount(ctx context.Context, id uuid.UUID) (int64, error)
}
