package repository

import (
	"context"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/characteristics"
	"github.com/google/uuid"
)

// ReadRepository описывает репозиторий чтения проекций (Read Model).
// Предоставляет доступ к денормализованным Read-моделям.
type ReadRepository[T any] interface {
	// GetByIDAsync получает Read-модель по идентификатору.
	GetByIDAsync(ctx context.Context, id uuid.UUID) (T, error)

	// GetAllAsync получает страницу Read-моделей.
	GetAllAsync(ctx context.Context, paging characteristics.Paging) ([]T, error)

	// CountAsync получает количество Read-моделей.
	CountAsync(ctx context.Context) (int64, error)
}
