// Package definition содержит определения ограниченных контекстов.
package definition

// BoundedContextDescription описывает ограниченный контекст.
type BoundedContextDescription struct {
	// ContextName — имя ограниченного контекста.
	ContextName string

	// MicroserviceVersion — версия микросервиса.
	MicroserviceVersion int
}

// NewBoundedContextDescription создаёт описание ограниченного контекста.
func NewBoundedContextDescription(contextName string, microserviceVersion int) BoundedContextDescription {
	return BoundedContextDescription{
		ContextName:         contextName,
		MicroserviceVersion: microserviceVersion,
	}
}
