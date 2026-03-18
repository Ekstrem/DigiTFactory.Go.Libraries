package characteristics

// CommandSubject описывает субъект изменений.
type CommandSubject interface {
	// CommandName возвращает имя метода агрегата, который вызывает команда.
	CommandName() string

	// SubjectName возвращает имя субъекта бизнес-операции.
	SubjectName() string
}
