package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FullEventSourcingRepository — стратегия Full Event Sourcing.
// Все события сохраняются, агрегат восстанавливается из полного стрима.
type FullEventSourcingRepository struct {
	pool   *pgxpool.Pool
	opts   Options
	logger *slog.Logger
}

// NewFullEventSourcingRepository создаёт репозиторий Full Event Sourcing.
func NewFullEventSourcingRepository(pool *pgxpool.Pool, opts Options, logger *slog.Logger) *FullEventSourcingRepository {
	return &FullEventSourcingRepository{pool: pool, opts: opts, logger: logger}
}

// GetByID получает все события агрегата, отсортированные по версии.
func (r *FullEventSourcingRepository) GetByID(ctx context.Context, id uuid.UUID) ([]DomainEventEntry, error) {
	sql := fmt.Sprintf(`SELECT id, version, correlation_token, bounded_context,
		command_name, subject_name, changed_value_objects, result, created_at
		FROM %s.domain_events WHERE id = $1 ORDER BY version`, r.opts.SchemaName)

	rows, err := r.pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("get events by id: %w", err)
	}
	defer rows.Close()

	var entries []DomainEventEntry
	for rows.Next() {
		var e DomainEventEntry
		if err := rows.Scan(&e.ID, &e.Version, &e.CorrelationToken, &e.BoundedContext,
			&e.CommandName, &e.SubjectName, &e.ChangedValueObjectsJSON, &e.Result, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		entries = append(entries, e)
	}

	r.logger.Debug("FullES: found events", "aggregateId", id, "count", len(entries))
	return entries, nil
}

// GetByIDAndVersion получает событие по идентификатору и версии.
func (r *FullEventSourcingRepository) GetByIDAndVersion(ctx context.Context, id uuid.UUID, version int64) (*DomainEventEntry, error) {
	sql := fmt.Sprintf(`SELECT id, version, correlation_token, bounded_context,
		command_name, subject_name, changed_value_objects, result, created_at
		FROM %s.domain_events WHERE id = $1 AND version = $2`, r.opts.SchemaName)

	var e DomainEventEntry
	err := r.pool.QueryRow(ctx, sql, id, version).Scan(
		&e.ID, &e.Version, &e.CorrelationToken, &e.BoundedContext,
		&e.CommandName, &e.SubjectName, &e.ChangedValueObjectsJSON, &e.Result, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get event by id and version: %w", err)
	}
	return &e, nil
}

// GetByCorrelationToken получает событие по маркеру корреляции.
func (r *FullEventSourcingRepository) GetByCorrelationToken(ctx context.Context, correlationToken uuid.UUID) (*DomainEventEntry, error) {
	sql := fmt.Sprintf(`SELECT id, version, correlation_token, bounded_context,
		command_name, subject_name, changed_value_objects, result, created_at
		FROM %s.domain_events WHERE correlation_token = $1 LIMIT 1`, r.opts.SchemaName)

	var e DomainEventEntry
	err := r.pool.QueryRow(ctx, sql, correlationToken).Scan(
		&e.ID, &e.Version, &e.CorrelationToken, &e.BoundedContext,
		&e.CommandName, &e.SubjectName, &e.ChangedValueObjectsJSON, &e.Result, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get event by correlation token: %w", err)
	}
	return &e, nil
}

// SaveEvent сохраняет событие в Event Store.
func (r *FullEventSourcingRepository) SaveEvent(ctx context.Context, entry *DomainEventEntry) error {
	sql := fmt.Sprintf(`INSERT INTO %s.domain_events
		(id, version, correlation_token, bounded_context, command_name, subject_name,
		 changed_value_objects, result, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9)`, r.opts.SchemaName)

	_, err := r.pool.Exec(ctx, sql,
		entry.ID, entry.Version, entry.CorrelationToken, entry.BoundedContext,
		entry.CommandName, entry.SubjectName, entry.ChangedValueObjectsJSON,
		entry.Result, entry.CreatedAt)
	if err != nil {
		return fmt.Errorf("save event: %w", err)
	}

	r.logger.Info("FullES: event saved", "aggregateId", entry.ID, "version", entry.Version)
	return nil
}

// SaveSnapshot — не используется в стратегии FullEventSourcing.
func (r *FullEventSourcingRepository) SaveSnapshot(_ context.Context, _ uuid.UUID, _ int64, _ string) error {
	r.logger.Warn("FullES: SaveSnapshot called but not used in FullEventSourcing strategy")
	return nil
}

// GetEventCount возвращает количество событий для агрегата.
func (r *FullEventSourcingRepository) GetEventCount(ctx context.Context, id uuid.UUID) (int, error) {
	sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s.domain_events WHERE id = $1`, r.opts.SchemaName)
	var count int
	err := r.pool.QueryRow(ctx, sql, id).Scan(&count)
	return count, err
}
