package scylla

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// ReadRepository — ScyllaDB/Cassandra реализация репозитория чтения Read-моделей.
type ReadRepository[T any] struct {
	session   *gocql.Session
	opts      Options
	modelType string
	logger    *slog.Logger
}

// NewReadRepository создаёт ScyllaDB репозиторий чтения.
func NewReadRepository[T any](session *gocql.Session, opts Options, modelType string, logger *slog.Logger) *ReadRepository[T] {
	return &ReadRepository[T]{session: session, opts: opts, modelType: modelType, logger: logger}
}

// GetByIDAsync получает Read-модель по идентификатору.
func (r *ReadRepository[T]) GetByIDAsync(id uuid.UUID) (T, error) {
	var result T
	cql := fmt.Sprintf(`SELECT data FROM %s.%s WHERE id = ?`, r.opts.Keyspace, r.opts.TableName)

	var data string
	if err := r.session.Query(cql, id).Scan(&data); err != nil {
		return result, fmt.Errorf("scylla get by id: %w", err)
	}

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return result, fmt.Errorf("unmarshal read model: %w", err)
	}
	return result, nil
}

// GetAllAsync получает все Read-модели данного типа.
func (r *ReadRepository[T]) GetAllAsync(pageSize int) ([]T, error) {
	cql := fmt.Sprintf(`SELECT data FROM %s.%s WHERE model_type = ? ALLOW FILTERING`,
		r.opts.Keyspace, r.opts.TableName)

	iter := r.session.Query(cql, r.modelType).PageSize(pageSize).Iter()

	var results []T
	var data string
	for iter.Scan(&data) {
		var item T
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("scylla get all: %w", err)
	}
	return results, nil
}

// CountAsync возвращает количество Read-моделей данного типа.
func (r *ReadRepository[T]) CountAsync() (int64, error) {
	cql := fmt.Sprintf(`SELECT COUNT(*) FROM %s.%s WHERE model_type = ? ALLOW FILTERING`,
		r.opts.Keyspace, r.opts.TableName)

	var count int64
	if err := r.session.Query(cql, r.modelType).Scan(&count); err != nil {
		return 0, fmt.Errorf("scylla count: %w", err)
	}
	return count, nil
}
