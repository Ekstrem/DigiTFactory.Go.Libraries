package characteristics

import "github.com/google/uuid"

// ComplexKey описывает комплексный ключ агрегата с типами по умолчанию (uuid.UUID, int64).
// Объединяет идентификатор, версию и маркер корреляции.
type ComplexKey interface {
	HasGuidKey
	HasInt64Version
	HasGuidCorrelationToken
}

// ComplexKeyGeneric описывает комплексный ключ агрегата с произвольными типами.
type ComplexKeyGeneric[TKey comparable, TVersion any] interface {
	HasKey[TKey]
	HasVersion[TVersion]
}

// complexKeyImpl — реализация комплексного ключа.
type complexKeyImpl struct {
	id               uuid.UUID
	version          int64
	correlationToken uuid.UUID
}

// NewComplexKey создаёт экземпляр комплексного ключа.
func NewComplexKey(id uuid.UUID, version int64, correlationToken uuid.UUID) ComplexKey {
	return &complexKeyImpl{
		id:               id,
		version:          version,
		correlationToken: correlationToken,
	}
}

// NewComplexKeyFromCorrelation создаёт экземпляр комплексного ключа,
// используя маркер корреляции в качестве идентификатора.
func NewComplexKeyFromCorrelation(correlationToken uuid.UUID, version int64) ComplexKey {
	return &complexKeyImpl{
		id:               correlationToken,
		version:          version,
		correlationToken: correlationToken,
	}
}

// ID возвращает идентификатор сущности.
func (k *complexKeyImpl) ID() uuid.UUID { return k.id }

// Version возвращает версию агрегата (Unix-миллисекунды).
func (k *complexKeyImpl) Version() int64 { return k.version }

// CorrelationToken возвращает идентификатор команды, создавшей новую версию.
func (k *complexKeyImpl) CorrelationToken() uuid.UUID { return k.correlationToken }
