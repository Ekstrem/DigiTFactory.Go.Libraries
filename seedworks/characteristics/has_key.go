// Package characteristics определяет базовые интерфейсы-характеристики доменных объектов.
package characteristics

import "github.com/google/uuid"

// HasKey указывает, что объект имеет идентификационное поле типа K.
type HasKey[K comparable] interface {
	// ID возвращает идентификатор объекта.
	ID() K
}

// HasGuidKey указывает, что объект имеет идентификационное поле типа uuid.UUID.
type HasGuidKey interface {
	HasKey[uuid.UUID]
}
