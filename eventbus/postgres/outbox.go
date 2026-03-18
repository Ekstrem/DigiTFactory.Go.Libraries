package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitializeOutbox создаёт таблицу outbox, если AutoCreateTable = true.
func InitializeOutbox(ctx context.Context, pool *pgxpool.Pool, opts Options, logger *slog.Logger) error {
	if !opts.AutoCreateTable {
		return nil
	}

	sql := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s.%s (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bounded_context TEXT NOT NULL,
    aggregate_id UUID NOT NULL,
    version BIGINT NOT NULL,
    correlation_token UUID NOT NULL,
    command_name TEXT NOT NULL,
    subject_name TEXT NOT NULL,
    command_version BIGINT NOT NULL,
    changed_value_objects JSONB NOT NULL DEFAULT '{}',
    result TEXT NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL
);
CREATE INDEX IF NOT EXISTS ix_%s_unprocessed
    ON %s.%s (created_at) WHERE processed_at IS NULL;
`, opts.SchemaName, opts.TableName, opts.TableName, opts.SchemaName, opts.TableName)

	_, err := pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("initialize outbox table: %w", err)
	}

	logger.Info("outbox table initialized",
		"schema", opts.SchemaName,
		"table", opts.TableName)
	return nil
}
