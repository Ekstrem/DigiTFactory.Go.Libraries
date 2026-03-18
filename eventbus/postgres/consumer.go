package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Consumer — polling consumer для outbox таблицы PostgreSQL.
type Consumer struct {
	mu       sync.RWMutex
	handlers map[string][]events.DomainEventHandler
	pool     *pgxpool.Pool
	options  Options
	logger   *slog.Logger
	cancel   context.CancelFunc
}

// NewConsumer создаёт Postgres consumer.
func NewConsumer(pool *pgxpool.Pool, options Options, logger *slog.Logger) *Consumer {
	return &Consumer{
		handlers: make(map[string][]events.DomainEventHandler),
		pool:     pool,
		options:  options,
		logger:   logger,
	}
}

// Subscribe подписывает обработчик на события ограниченного контекста.
func (c *Consumer) Subscribe(boundedContext string, handler events.DomainEventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[boundedContext] = append(c.handlers[boundedContext], handler)
}

// Unsubscribe отписывает обработчик.
func (c *Consumer) Unsubscribe(boundedContext string, handler events.DomainEventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	handlers := c.handlers[boundedContext]
	for i, h := range handlers {
		if h == handler {
			c.handlers[boundedContext] = append(handlers[:i], handlers[i+1:]...)
			return
		}
	}
}

// StartConsuming запускает polling loop.
func (c *Consumer) StartConsuming(ctx context.Context) {
	childCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	go c.pollLoop(childCtx)
	c.logger.Info("started postgres outbox consumer")
}

// Stop останавливает polling loop.
func (c *Consumer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Consumer) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(c.options.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.pollBatch(ctx); err != nil {
				c.logger.Error("outbox poll error", "error", err)
			}
		}
	}
}

func (c *Consumer) pollBatch(ctx context.Context) error {
	sql := fmt.Sprintf(`
SELECT id, bounded_context, aggregate_id, version, correlation_token,
       command_name, subject_name, command_version,
       changed_value_objects, result, reason, created_at
FROM %s.%s
WHERE processed_at IS NULL
ORDER BY created_at
LIMIT %d`,
		c.options.SchemaName, c.options.TableName, c.options.BatchSize)

	rows, err := c.pool.Query(ctx, sql)
	if err != nil {
		return fmt.Errorf("query outbox: %w", err)
	}
	defer rows.Close()

	var processedIDs []any
	for rows.Next() {
		var env Envelope
		var id string
		if err := rows.Scan(
			&id, &env.BoundedContext, &env.AggregateID, &env.Version,
			&env.CorrelationToken, &env.CommandName, &env.SubjectName,
			&env.CommandVersion, &env.ChangedValueObjectsJSON,
			&env.Result, &env.Reason, &env.CreatedAt,
		); err != nil {
			c.logger.Error("scan outbox row", "error", err)
			continue
		}

		c.dispatchEnvelope(ctx, &env)
		processedIDs = append(processedIDs, id)
	}

	if len(processedIDs) > 0 {
		return c.markProcessed(ctx, processedIDs)
	}
	return nil
}

func (c *Consumer) dispatchEnvelope(ctx context.Context, env *Envelope) {
	c.mu.RLock()
	handlers := c.handlers[env.BoundedContext]
	snapshot := make([]events.DomainEventHandler, len(handlers))
	copy(snapshot, handlers)
	c.mu.RUnlock()

	var changedVOs map[string]any
	_ = json.Unmarshal([]byte(env.ChangedValueObjectsJSON), &changedVOs)

	event := &events.DomainEvent{
		AggregateID:         env.AggregateID,
		Ver:                 env.Version,
		BoundedContextName:  env.BoundedContext,
		ContextName:         env.BoundedContext,
		ChangedValueObjects: changedVOs,
		ResultReason:        env.Reason,
		Command: events.NewCommandToAggregate(
			env.CorrelationToken, env.CommandName, env.SubjectName, env.CommandVersion,
		),
	}

	for _, handler := range snapshot {
		if err := handler.HandleEvent(ctx, event); err != nil {
			c.logger.Error("dispatch outbox event error", "error", err)
		}
	}
}

func (c *Consumer) markProcessed(ctx context.Context, ids []any) error {
	sql := fmt.Sprintf(`
UPDATE %s.%s SET processed_at = NOW() WHERE id = ANY($1)`,
		c.options.SchemaName, c.options.TableName)

	_, err := c.pool.Exec(ctx, sql, ids)
	return err
}
