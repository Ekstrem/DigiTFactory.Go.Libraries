// Package mongo предоставляет MongoDB реализацию Event Store.
package mongo

import "github.com/Ekstrem/DigiTFactory.Go.Libraries/commandrepository"

// Options — настройки MongoDB Event Store.
type Options struct {
	// ConnectionString — строка подключения к MongoDB.
	ConnectionString string

	// DatabaseName — имя базы данных.
	DatabaseName string

	// Strategy — стратегия хранения событий.
	Strategy commandrepository.EventStoreStrategy

	// SnapshotInterval — интервал создания snapshot (только для SnapshotAfterN).
	SnapshotInterval int

	// EventsCollection — имя коллекции событий.
	EventsCollection string

	// SnapshotsCollection — имя коллекции snapshot'ов.
	SnapshotsCollection string

	// StatesCollection — имя коллекции состояний агрегатов.
	StatesCollection string
}

// DefaultOptions возвращает настройки по умолчанию.
func DefaultOptions() Options {
	return Options{
		ConnectionString:    "mongodb://localhost:27017",
		DatabaseName:        "event_store",
		Strategy:            commandrepository.FullEventSourcing,
		SnapshotInterval:    10,
		EventsCollection:    "domain_events",
		SnapshotsCollection: "snapshots",
		StatesCollection:    "aggregate_states",
	}
}
