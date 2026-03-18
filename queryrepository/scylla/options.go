// Package scylla предоставляет ScyllaDB/Cassandra реализацию Read Store.
package scylla

// Options — настройки ScyllaDB Read Store.
type Options struct {
	// ContactPoints — точки подключения к ScyllaDB/Cassandra.
	ContactPoints []string

	// Port — порт ScyllaDB/Cassandra.
	Port int

	// Keyspace — имя keyspace.
	Keyspace string

	// TableName — имя таблицы для Read-моделей.
	TableName string

	// ReplicationStrategy — стратегия репликации.
	ReplicationStrategy string

	// ReplicationFactor — фактор репликации.
	ReplicationFactor int

	// AutoCreateSchema — автоматически создавать keyspace и таблицу при старте.
	AutoCreateSchema bool
}

// DefaultOptions возвращает настройки по умолчанию.
func DefaultOptions() Options {
	return Options{
		ContactPoints:       []string{"localhost"},
		Port:                9042,
		Keyspace:            "read_models",
		TableName:           "projections",
		ReplicationStrategy: "SimpleStrategy",
		ReplicationFactor:   1,
		AutoCreateSchema:    true,
	}
}
