package tacticalpatterns

import (
	"context"

	"github.com/google/uuid"
)

// AggregateProvider описывает провайдер получения экземпляра агрегата.
type AggregateProvider interface {
	// GetAggregate получает агрегат по идентификатору и версии.
	GetAggregate(id uuid.UUID, version int64) (AnemicModel, error)

	// GetAggregateAsync получает агрегат по идентификатору и версии асинхронно.
	GetAggregateAsync(ctx context.Context, id uuid.UUID, version int64) (AnemicModel, error)

	// GetAggregateByID получает агрегат по идентификатору.
	GetAggregateByID(id uuid.UUID) (AnemicModel, error)

	// GetAggregateByIDAsync получает агрегат по идентификатору асинхронно.
	GetAggregateByIDAsync(ctx context.Context, id uuid.UUID) (AnemicModel, error)
}
