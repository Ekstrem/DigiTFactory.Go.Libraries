package characteristics

// CommandToAggregate описывает метаданные команды к агрегату.
// Объединяет версионность, маркер корреляции и описание субъекта команды.
type CommandToAggregate interface {
	HasInt64Version
	HasGuidCorrelationToken
	CommandSubject
}
