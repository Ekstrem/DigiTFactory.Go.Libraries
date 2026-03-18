// Package tacticalpatterns содержит тактические паттерны DDD:
// объект-значение, сущность, агрегат, анемичная модель и бизнес-операции.
package tacticalpatterns

// ValueObject — маркерный интерфейс объекта-значения.
// В Go реализуется как пустой интерфейс; типы, представляющие объект-значения,
// должны документировать это соглашение.
type ValueObject = any

// ValueObjectWrapper — обёртка над простыми типами,
// чтобы предоставить их как объект-значение.
type ValueObjectWrapper[T any] struct {
	// Value — значение, которое необходимо предоставить как объект-значение.
	Value T
}

// NewValueObjectWrapper создаёт обёртку над значением.
func NewValueObjectWrapper[T any](value T) ValueObjectWrapper[T] {
	return ValueObjectWrapper[T]{Value: value}
}
