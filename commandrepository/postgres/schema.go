package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitializeSchema создаёт таблицы Event Store в PostgreSQL.
func InitializeSchema(ctx context.Context, pool *pgxpool.Pool, opts Options, logger *slog.Logger) error {
	schema := opts.SchemaName

	sql := fmt.Sprintf(`
CREATE SCHEMA IF NOT EXISTS %s;

CREATE TABLE IF NOT EXISTS %s.domain_events (
    id UUID NOT NULL,
    version BIGINT NOT NULL,
    correlation_token UUID NOT NULL,
    bounded_context TEXT NOT NULL,
    command_name TEXT NOT NULL,
    subject_name TEXT NOT NULL,
    changed_value_objects JSONB NOT NULL DEFAULT '{}',
    result TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, version)
);

CREATE INDEX IF NOT EXISTS ix_domain_events_correlation
    ON %s.domain_events (correlation_token);

CREATE TABLE IF NOT EXISTS %s.snapshots (
    id UUID NOT NULL,
    version BIGINT NOT NULL,
    aggregate_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, version)
);

CREATE TABLE IF NOT EXISTS %s.aggregate_states (
    id UUID PRIMARY KEY,
    version BIGINT NOT NULL,
    aggregate_json JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`, schema, schema, schema, schema, schema)

	_, err := pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("initialize event store schema: %w", err)
	}

	logger.Info("event store schema initialized", "schema", schema)
	return nil
}
