package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StateOnlyRepository — стратегия StateOnly.
// Каждый раз сохраняется агрегат целиком (без истории событий).
type StateOnlyRepository struct {
	pool   *pgxpool.Pool
	opts   Options
	logger *slog.Logger
}

// NewStateOnlyRepository создаёт репозиторий со стратегией StateOnly.
func NewStateOnlyRepository(pool *pgxpool.Pool, opts Options, logger *slog.Logger) *StateOnlyRepository {
	return &StateOnlyRepository{pool: pool, opts: opts, logger: logger}
}

// GetByID получает текущее состояние агрегата (вернёт один элемент).
func (r *StateOnlyRepository) GetByID(ctx context.Context, id uuid.UUID) ([]DomainEventEntry, error) {
	// StateOnly не хранит события, а только текущее состояние
	return nil, nil
}

// GetByIDAndVersion не поддерживается в стратегии StateOnly.
func (r *StateOnlyRepository) GetByIDAndVersion(_ context.Context, _ uuid.UUID, _ int64) (*DomainEventEntry, error) {
	return nil, fmt.Errorf("GetByIDAndVersion not supported in StateOnly strategy")
}

// GetByCorrelationToken не поддерживается в стратегии StateOnly.
func (r *StateOnlyRepository) GetByCorrelationToken(_ context.Context, _ uuid.UUID) (*DomainEventEntry, error) {
	return nil, fmt.Errorf("GetByCorrelationToken not supported in StateOnly strategy")
}

// SaveEvent в стратегии StateOnly сохраняет/обновляет текущее состояние агрегата.
func (r *StateOnlyRepository) SaveEvent(ctx context.Context, entry *DomainEventEntry) error {
	sql := fmt.Sprintf(`INSERT INTO %s.aggregate_states (id, version, aggregate_json, updated_at)
		VALUES ($1, $2, $3::jsonb, $4)
		ON CONFLICT (id) DO UPDATE SET
			version = EXCLUDED.version,
			aggregate_json = EXCLUDED.aggregate_json,
			updated_at = EXCLUDED.updated_at`, r.opts.SchemaName)

	_, err := r.pool.Exec(ctx, sql,
		entry.ID, entry.Version, entry.ChangedValueObjectsJSON, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("save aggregate state: %w", err)
	}

	r.logger.Info("StateOnly: state saved", "aggregateId", entry.ID, "version", entry.Version)
	return nil
}

// SaveSnapshot — не используется в стратегии StateOnly.
func (r *StateOnlyRepository) SaveSnapshot(_ context.Context, _ uuid.UUID, _ int64, _ string) error {
	return nil
}

// GetEventCount — всегда 1 в стратегии StateOnly (есть или нет состояние).
func (r *StateOnlyRepository) GetEventCount(ctx context.Context, id uuid.UUID) (int, error) {
	sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s.aggregate_states WHERE id = $1`, r.opts.SchemaName)
	var count int
	err := r.pool.QueryRow(ctx, sql, id).Scan(&count)
	return count, err
}
