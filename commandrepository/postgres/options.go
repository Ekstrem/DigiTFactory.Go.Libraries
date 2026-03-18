// Package postgres предоставляет PostgreSQL реализацию Event Store.
package postgres

import "github.com/Ekstrem/DigiTFactory.Go.Libraries/commandrepository"

// Options — настройки PostgreSQL Event Store.
type Options struct {
	// Strategy — стратегия хранения событий.
	Strategy commandrepository.EventStoreStrategy

	// SnapshotInterval — интервал создания snapshot (только для SnapshotAfterN).
	SnapshotInterval int

	// SchemaName — имя схемы PostgreSQL для таблиц Event Store.
	SchemaName string
}

// DefaultOptions возвращает настройки по умолчанию.
func DefaultOptions() Options {
	return Options{
		Strategy:         commandrepository.FullEventSourcing,
		SnapshotInterval: 10,
		SchemaName:       "commands",
	}
}
