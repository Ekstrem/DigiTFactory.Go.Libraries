package scylla

import (
	"fmt"
	"log/slog"

	"github.com/gocql/gocql"
)

// InitializeSchema создаёт keyspace и таблицу в ScyllaDB, если AutoCreateSchema = true.
func InitializeSchema(session *gocql.Session, opts Options, logger *slog.Logger) error {
	if !opts.AutoCreateSchema {
		return nil
	}

	createKeyspace := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH replication = {
			'class': '%s',
			'replication_factor': %d
		}`, opts.Keyspace, opts.ReplicationStrategy, opts.ReplicationFactor)

	if err := session.Query(createKeyspace).Exec(); err != nil {
		return fmt.Errorf("create keyspace: %w", err)
	}

	createTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			id uuid PRIMARY KEY,
			data text,
			model_type text,
			updated_at timestamp
		)`, opts.Keyspace, opts.TableName)

	if err := session.Query(createTable).Exec(); err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	logger.Info("ScyllaDB schema initialized",
		"keyspace", opts.Keyspace,
		"table", opts.TableName)
	return nil
}
