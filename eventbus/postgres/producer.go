package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Producer — INSERT доменных событий в outbox таблицу PostgreSQL.
type Producer struct {
	pool    *pgxpool.Pool
	options Options
	logger  *slog.Logger
}

// NewProducer создаёт Postgres producer.
func NewProducer(pool *pgxpool.Pool, options Options, logger *slog.Logger) *Producer {
	return &Producer{
		pool:    pool,
		options: options,
		logger:  logger,
	}
}

// Publish публикует доменное событие в outbox таблицу.
func (p *Producer) Publish(ctx context.Context, event *events.DomainEvent) error {
	envelope, err := ToEnvelope(event)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf(`
INSERT INTO %s.%s
    (bounded_context, aggregate_id, version, correlation_token,
     command_name, subject_name, command_version,
     changed_value_objects, result, reason, created_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9, $10, $11)`,
		p.options.SchemaName, p.options.TableName)

	_, err = p.pool.Exec(ctx, sql,
		envelope.BoundedContext,
		envelope.AggregateID,
		envelope.Version,
		envelope.CorrelationToken,
		envelope.CommandName,
		envelope.SubjectName,
		envelope.CommandVersion,
		envelope.ChangedValueObjectsJSON,
		envelope.Result,
		envelope.Reason,
		envelope.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert into outbox: %w", err)
	}

	p.logger.Debug("published event to outbox",
		"schema", p.options.SchemaName,
		"table", p.options.TableName,
		"aggregateId", envelope.AggregateID)
	return nil
}
