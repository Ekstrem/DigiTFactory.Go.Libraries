package characteristics

import "github.com/google/uuid"

// HasCorrelationToken указывает, что объект содержит маркер корреляции типа T.
type HasCorrelationToken[T any] interface {
	// CorrelationToken возвращает идентификатор команды, создавшей новую версию.
	CorrelationToken() T
}

// HasGuidCorrelationToken указывает, что объект содержит маркер корреляции типа uuid.UUID.
type HasGuidCorrelationToken interface {
	HasCorrelationToken[uuid.UUID]
}
