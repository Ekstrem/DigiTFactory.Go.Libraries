package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReadModelStore — PostgreSQL реализация хранилища Read-моделей.
type ReadModelStore struct {
	pool      *pgxpool.Pool
	tableName string
	logger    *slog.Logger
}

// NewReadModelStore создаёт PostgreSQL хранилище Read-моделей.
func NewReadModelStore(pool *pgxpool.Pool, tableName string, logger *slog.Logger) *ReadModelStore {
	return &ReadModelStore{pool: pool, tableName: tableName, logger: logger}
}

// UpsertAsync вставляет или обновляет Read-модель.
func (s *ReadModelStore) UpsertAsync(ctx context.Context, id uuid.UUID, model any) error {
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("marshal read model: %w", err)
	}

	sql := fmt.Sprintf(`INSERT INTO %s (id, data, updated_at)
		VALUES ($1, $2::jsonb, NOW())
		ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data, updated_at = NOW()`, s.tableName)

	_, err = s.pool.Exec(ctx, sql, id, string(data))
	if err != nil {
		return fmt.Errorf("upsert read model: %w", err)
	}

	s.logger.Debug("upserted read model", "table", s.tableName, "id", id)
	return nil
}

// DeleteAsync удаляет Read-модель по идентификатору.
func (s *ReadModelStore) DeleteAsync(ctx context.Context, id uuid.UUID) error {
	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, s.tableName)
	_, err := s.pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete read model: %w", err)
	}

	s.logger.Debug("deleted read model", "table", s.tableName, "id", id)
	return nil
}
