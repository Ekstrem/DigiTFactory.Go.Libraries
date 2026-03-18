// Package redis предоставляет Redis реализацию Read Store.
package redis

// Options — настройки Redis Read Store.
type Options struct {
	// ConnectionString — строка подключения к Redis.
	ConnectionString string

	// KeyPrefix — префикс ключей в Redis.
	KeyPrefix string
}

// DefaultOptions возвращает настройки по умолчанию.
func DefaultOptions() Options {
	return Options{
		ConnectionString: "localhost:6379",
		KeyPrefix:        "readmodel:",
	}
}
