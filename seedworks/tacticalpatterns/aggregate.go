package tacticalpatterns

// Aggregate — агрегат, предоставляющий доступ к бизнес-операциям ограниченного контекста.
type Aggregate struct {
	// Operations — бизнес-операции агрегата.
	Operations map[string]AggregateBusinessOperation
}

// NewAggregate создаёт экземпляр агрегата из области ограниченного контекста.
func NewAggregate(scope BoundedContextScope) *Aggregate {
	return &Aggregate{
		Operations: scope.GetOperations(),
	}
}
