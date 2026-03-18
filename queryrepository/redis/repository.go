package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

// ReadRepository — Redis реализация репозитория чтения Read-моделей.
type ReadRepository[T any] struct {
	client    *goredis.Client
	keyPrefix string
	logger    *slog.Logger
}

// NewReadRepository создаёт Redis репозиторий чтения.
func NewReadRepository[T any](client *goredis.Client, keyPrefix string, logger *slog.Logger) *ReadRepository[T] {
	return &ReadRepository[T]{client: client, keyPrefix: keyPrefix, logger: logger}
}

// GetByIDAsync получает Read-модель по идентификатору.
func (r *ReadRepository[T]) GetByIDAsync(ctx context.Context, id uuid.UUID) (T, error) {
	var result T
	key := r.keyPrefix + id.String()

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return result, fmt.Errorf("redis get: %w", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("unmarshal read model: %w", err)
	}
	return result, nil
}

// CountAsync возвращает приблизительное количество Read-моделей по шаблону ключа.
func (r *ReadRepository[T]) CountAsync(ctx context.Context) (int64, error) {
	pattern := r.keyPrefix + "*"
	var count int64
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return 0, fmt.Errorf("redis scan: %w", err)
		}
		count += int64(len(keys))
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return count, nil
}
