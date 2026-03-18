package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SnapshotRepository — стратегия SnapshotAfterN.
// Каждые N событий сохраняется snapshot, при чтении: snapshot + события после.
type SnapshotRepository struct {
	pool   *pgxpool.Pool
	opts   Options
	logger *slog.Logger
}

// NewSnapshotRepository создаёт репозиторий со стратегией SnapshotAfterN.
func NewSnapshotRepository(pool *pgxpool.Pool, opts Options, logger *slog.Logger) *SnapshotRepository {
	return &SnapshotRepository{pool: pool, opts: opts, logger: logger}
}

// GetByID получает все события агрегата начиная с последнего snapshot.
func (r *SnapshotRepository) GetByID(ctx context.Context, id uuid.UUID) ([]DomainEventEntry, error) {
	// Находим последний snapshot
	snapshotSQL := fmt.Sprintf(`SELECT version FROM %s.snapshots
		WHERE id = $1 ORDER BY version DESC LIMIT 1`, r.opts.SchemaName)
	var snapshotVersion int64
	err := r.pool.QueryRow(ctx, snapshotSQL, id).Scan(&snapshotVersion)
	if err != nil {
		snapshotVersion = 0 // Нет snapshot, читаем с начала
	}

	sql := fmt.Sprintf(`SELECT id, version, correlation_token, bounded_context,
		command_name, subject_name, changed_value_objects, result, created_at
		FROM %s.domain_events WHERE id = $1 AND version > $2 ORDER BY version`, r.opts.SchemaName)

	rows, err := r.pool.Query(ctx, sql, id, snapshotVersion)
	if err != nil {
		return nil, fmt.Errorf("get events after snapshot: %w", err)
	}
	defer rows.Close()

	var entries []DomainEventEntry
	for rows.Next() {
		var e DomainEventEntry
		if err := rows.Scan(&e.ID, &e.Version, &e.CorrelationToken, &e.BoundedContext,
			&e.CommandName, &e.SubjectName, &e.ChangedValueObjectsJSON, &e.Result, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// GetByIDAndVersion получает событие по идентификатору и версии.
func (r *SnapshotRepository) GetByIDAndVersion(ctx context.Context, id uuid.UUID, version int64) (*DomainEventEntry, error) {
	sql := fmt.Sprintf(`SELECT id, version, correlation_token, bounded_context,
		command_name, subject_name, changed_value_objects, result, created_at
		FROM %s.domain_events WHERE id = $1 AND version = $2`, r.opts.SchemaName)

	var e DomainEventEntry
	err := r.pool.QueryRow(ctx, sql, id, version).Scan(
		&e.ID, &e.Version, &e.CorrelationToken, &e.BoundedContext,
		&e.CommandName, &e.SubjectName, &e.ChangedValueObjectsJSON, &e.Result, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetByCorrelationToken получает событие по маркеру корреляции.
func (r *SnapshotRepository) GetByCorrelationToken(ctx context.Context, ct uuid.UUID) (*DomainEventEntry, error) {
	sql := fmt.Sprintf(`SELECT id, version, correlation_token, bounded_context,
		command_name, subject_name, changed_value_objects, result, created_at
		FROM %s.domain_events WHERE correlation_token = $1 LIMIT 1`, r.opts.SchemaName)

	var e DomainEventEntry
	err := r.pool.QueryRow(ctx, sql, ct).Scan(
		&e.ID, &e.Version, &e.CorrelationToken, &e.BoundedContext,
		&e.CommandName, &e.SubjectName, &e.ChangedValueObjectsJSON, &e.Result, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// SaveEvent сохраняет событие и при необходимости создаёт snapshot.
func (r *SnapshotRepository) SaveEvent(ctx context.Context, entry *DomainEventEntry) error {
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

	r.logger.Info("Snapshot: event saved", "aggregateId", entry.ID, "version", entry.Version)
	return nil
}

// SaveSnapshot сохраняет snapshot агрегата.
func (r *SnapshotRepository) SaveSnapshot(ctx context.Context, id uuid.UUID, version int64, aggregateJSON string) error {
	sql := fmt.Sprintf(`INSERT INTO %s.snapshots (id, version, aggregate_json, created_at)
		VALUES ($1, $2, $3::jsonb, $4)
		ON CONFLICT (id, version) DO UPDATE SET aggregate_json = EXCLUDED.aggregate_json`, r.opts.SchemaName)

	_, err := r.pool.Exec(ctx, sql, id, version, aggregateJSON, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}
	r.logger.Info("Snapshot saved", "aggregateId", id, "version", version)
	return nil
}

// GetEventCount возвращает количество событий для агрегата.
func (r *SnapshotRepository) GetEventCount(ctx context.Context, id uuid.UUID) (int, error) {
	sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s.domain_events WHERE id = $1`, r.opts.SchemaName)
	var count int
	err := r.pool.QueryRow(ctx, sql, id).Scan(&count)
	return count, err
}
