// Package kafka предоставляет Apache Kafka реализацию EventBus.
package kafka

// Options — настройки Kafka EventBus.
type Options struct {
	// BootstrapServers — адреса Kafka брокеров.
	BootstrapServers string

	// TopicPrefix — префикс имён топиков (например, "domain-events" → "domain-events.mycontext").
	TopicPrefix string

	// GroupID — идентификатор consumer group.
	GroupID string

	// AutoOffsetReset — позиция чтения при отсутствии offset ("earliest" или "latest").
	AutoOffsetReset string
}

// DefaultOptions возвращает настройки по умолчанию.
func DefaultOptions() Options {
	return Options{
		BootstrapServers: "localhost:9092",
		TopicPrefix:      "domain-events",
		GroupID:          "default-group",
		AutoOffsetReset:  "earliest",
	}
}
