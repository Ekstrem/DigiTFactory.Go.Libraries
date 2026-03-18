package repository

import (
	"context"

	"github.com/google/uuid"
)

// AnemicModelRepository описывает репозиторий получения из хранилища анемичной модели агрегата.
// В норме — зависимость для провайдера агрегатов (AggregateProvider).
type AnemicModelRepository[T any] interface {
	// GetByID получает стрим анемичных моделей по идентификатору.
	GetByID(ctx context.Context, id uuid.UUID) ([]T, error)

	// GetByIDAndVersion получает анемичную модель по идентификатору и версии.
	GetByIDAndVersion(ctx context.Context, id uuid.UUID, version int64) (T, error)

	// GetByCorrelationToken получает анемичную модель по маркеру корреляции.
	GetByCorrelationToken(ctx context.Context, correlationToken uuid.UUID) (T, error)
}
