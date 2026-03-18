package repository

import (
	"context"

	"github.com/google/uuid"
)

// ReadModelStore описывает хранилище записи Read-моделей (проекций).
// Используется проекциями для обновления денормализованных данных
// на основе доменных событий.
type ReadModelStore[T any] interface {
	// UpsertAsync вставляет или обновляет Read-модель.
	UpsertAsync(ctx context.Context, model T) error

	// DeleteAsync удаляет Read-модель по идентификатору.
	DeleteAsync(ctx context.Context, id uuid.UUID) error
}
