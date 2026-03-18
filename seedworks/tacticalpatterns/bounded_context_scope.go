package tacticalpatterns

// BoundedContextScope описывает область ограниченного контекста.
// Определяет набор бизнес-операций и валидаторов для агрегата.
type BoundedContextScope interface {
	// GetOperations возвращает бизнес-операции агрегата.
	GetOperations() map[string]AggregateBusinessOperation

	// GetValidators возвращает валидаторы бизнес-сущностей.
	GetValidators() []Validator
}
