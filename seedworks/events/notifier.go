package events

// DomainEventNotifier — базовая структура нотификатора о доменных событиях.
// Служит для проксирования вызовов к агрегату.
type DomainEventNotifier struct {
	// Aggregate — агрегат, вызовы к которому необходимо проксировать.
	Aggregate any
}

// NewDomainEventNotifier создаёт нотификатор доменных событий.
func NewDomainEventNotifier(aggregate any) *DomainEventNotifier {
	return &DomainEventNotifier{Aggregate: aggregate}
}
