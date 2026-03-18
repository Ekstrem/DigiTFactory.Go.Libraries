package scylla

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// ReadModelStore — ScyllaDB/Cassandra реализация хранилища Read-моделей.
type ReadModelStore struct {
	session   *gocql.Session
	opts      Options
	modelType string
	logger    *slog.Logger
}

// NewReadModelStore создаёт ScyllaDB хранилище Read-моделей.
func NewReadModelStore(session *gocql.Session, opts Options, modelType string, logger *slog.Logger) *ReadModelStore {
	return &ReadModelStore{session: session, opts: opts, modelType: modelType, logger: logger}
}

// UpsertAsync вставляет или обновляет Read-модель.
func (s *ReadModelStore) UpsertAsync(id uuid.UUID, model any) error {
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("marshal read model: %w", err)
	}

	cql := fmt.Sprintf(`INSERT INTO %s.%s (id, data, model_type, updated_at)
		VALUES (?, ?, ?, toTimestamp(now()))`, s.opts.Keyspace, s.opts.TableName)

	if err := s.session.Query(cql, id, string(data), s.modelType).Exec(); err != nil {
		return fmt.Errorf("scylla upsert: %w", err)
	}

	s.logger.Debug("upserted read model in scylla", "id", id, "modelType", s.modelType)
	return nil
}

// DeleteAsync удаляет Read-модель по идентификатору.
func (s *ReadModelStore) DeleteAsync(id uuid.UUID) error {
	cql := fmt.Sprintf(`DELETE FROM %s.%s WHERE id = ?`, s.opts.Keyspace, s.opts.TableName)

	if err := s.session.Query(cql, id).Exec(); err != nil {
		return fmt.Errorf("scylla delete: %w", err)
	}

	s.logger.Debug("deleted read model from scylla", "id", id)
	return nil
}
