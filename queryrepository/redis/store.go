package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

// ReadModelStore — Redis реализация хранилища Read-моделей.
type ReadModelStore struct {
	client    *goredis.Client
	keyPrefix string
	logger    *slog.Logger
}

// NewReadModelStore создаёт Redis хранилище Read-моделей.
func NewReadModelStore(client *goredis.Client, keyPrefix string, logger *slog.Logger) *ReadModelStore {
	return &ReadModelStore{client: client, keyPrefix: keyPrefix, logger: logger}
}

// UpsertAsync вставляет или обновляет Read-модель.
func (s *ReadModelStore) UpsertAsync(ctx context.Context, id uuid.UUID, model any) error {
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("marshal read model: %w", err)
	}

	key := s.keyPrefix + id.String()
	if err := s.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	s.logger.Debug("upserted read model in redis", "key", key)
	return nil
}

// DeleteAsync удаляет Read-модель по идентификатору.
func (s *ReadModelStore) DeleteAsync(ctx context.Context, id uuid.UUID) error {
	key := s.keyPrefix + id.String()
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}

	s.logger.Debug("deleted read model from redis", "key", key)
	return nil
}
