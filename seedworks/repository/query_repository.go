package repository

import (
	"context"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/characteristics"
)

// QueryRepository описывает репозиторий запросов (Read).
type QueryRepository[T any] interface {
	// Get возвращает страницу записей.
	Get(paging characteristics.Paging) ([]T, error)

	// Count возвращает количество записей.
	Count() (int64, error)

	// CountAsync возвращает количество записей асинхронно.
	CountAsync(ctx context.Context) (int64, error)
}
