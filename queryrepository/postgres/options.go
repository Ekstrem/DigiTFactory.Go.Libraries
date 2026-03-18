// Package postgres предоставляет PostgreSQL реализацию Read Store.
package postgres

// Options — настройки PostgreSQL Read Store.
type Options struct {
	// ConnectionString — строка подключения к PostgreSQL.
	ConnectionString string
}
