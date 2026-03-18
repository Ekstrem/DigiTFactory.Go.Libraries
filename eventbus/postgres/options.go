// Package postgres предоставляет PostgreSQL реализацию EventBus (outbox pattern).
package postgres

import "time"

// Options — настройки PostgreSQL EventBus (outbox pattern).
type Options struct {
	// ConnectionString — строка подключения к PostgreSQL.
	ConnectionString string

	// SchemaName — схема таблицы outbox.
	SchemaName string

	// TableName — имя таблицы outbox.
	TableName string

	// PollingInterval — интервал polling необработанных событий.
	PollingInterval time.Duration

	// BatchSize — максимальное количество событий за один poll.
	BatchSize int

	// AutoCreateTable — автоматически создавать таблицу при старте.
	AutoCreateTable bool
}

// DefaultOptions возвращает настройки по умолчанию.
func DefaultOptions() Options {
	return Options{
		SchemaName:      "public",
		TableName:       "domain_events_outbox",
		PollingInterval: 5 * time.Second,
		BatchSize:       100,
		AutoCreateTable: true,
	}
}
