package repository

import (
	"context"
	"io"
)

// IsolationLevel определяет уровень изоляции транзакции.
type IsolationLevel int

const (
	// ReadCommitted — чтение зафиксированных данных.
	ReadCommitted IsolationLevel = iota

	// ReadUncommitted — чтение незафиксированных данных.
	ReadUncommitted

	// RepeatableRead — повторяемое чтение.
	RepeatableRead

	// Serializable — сериализуемость.
	Serializable
)

// Transaction описывает управление транзакциями.
type Transaction interface {
	// Begin начинает транзакцию и возвращает объект для её завершения.
	Begin() (io.Closer, error)

	// BeginAsync начинает транзакцию асинхронно.
	BeginAsync(ctx context.Context) (io.Closer, error)

	// Commit подтверждает транзакцию.
	Commit() error

	// RollBack откатывает транзакцию.
	RollBack() error

	// SetIsolationLevel устанавливает уровень изоляции.
	SetIsolationLevel(level IsolationLevel)
}
