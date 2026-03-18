package tacticalpatterns

import "github.com/google/uuid"

// Entity описывает сущность с идентификатором типа uuid.UUID.
type Entity interface {
	// ID возвращает идентификатор сущности.
	ID() uuid.UUID
}

// EntityGeneric описывает сущность с идентификатором произвольного типа.
type EntityGeneric[TKey comparable] interface {
	// ID возвращает идентификатор сущности.
	ID() TKey
}
