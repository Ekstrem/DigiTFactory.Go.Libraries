package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/characteristics"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReadRepository — PostgreSQL реализация репозитория чтения Read-моделей.
type ReadRepository[T any] struct {
	pool      *pgxpool.Pool
	tableName string
	logger    *slog.Logger
}

// NewReadRepository создаёт PostgreSQL репозиторий чтения.
func NewReadRepository[T any](pool *pgxpool.Pool, tableName string, logger *slog.Logger) *ReadRepository[T] {
	return &ReadRepository[T]{pool: pool, tableName: tableName, logger: logger}
}

// GetByIDAsync получает Read-модель по идентификатору.
func (r *ReadRepository[T]) GetByIDAsync(ctx context.Context, id uuid.UUID) (T, error) {
	var result T
	sql := fmt.Sprintf(`SELECT data FROM %s WHERE id = $1`, r.tableName)

	var data string
	err := r.pool.QueryRow(ctx, sql, id).Scan(&data)
	if err != nil {
		return result, fmt.Errorf("get read model by id: %w", err)
	}

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return result, fmt.Errorf("unmarshal read model: %w", err)
	}
	return result, nil
}

// GetAllAsync получает страницу Read-моделей.
func (r *ReadRepository[T]) GetAllAsync(ctx context.Context, paging characteristics.Paging) ([]T, error) {
	offset := (paging.Page() - 1) * paging.PageSize()
	sql := fmt.Sprintf(`SELECT data FROM %s ORDER BY id LIMIT $1 OFFSET $2`, r.tableName)

	rows, err := r.pool.Query(ctx, sql, paging.PageSize(), offset)
	if err != nil {
		return nil, fmt.Errorf("get all read models: %w", err)
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		var item T
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

// CountAsync получает количество Read-моделей.
func (r *ReadRepository[T]) CountAsync(ctx context.Context) (int64, error) {
	sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, r.tableName)
	var count int64
	err := r.pool.QueryRow(ctx, sql).Scan(&count)
	return count, err
}
