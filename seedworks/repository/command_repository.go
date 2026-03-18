// Package repository содержит интерфейсы репозиториев CQRS.
package repository

// CommandRepository описывает репозиторий команд (CUD).
type CommandRepository[T any] interface {
	// Add добавляет запись.
	Add(entity T)

	// Update обновляет запись.
	Update(entity T)

	// Delete удаляет запись.
	Delete(entity T)
}
