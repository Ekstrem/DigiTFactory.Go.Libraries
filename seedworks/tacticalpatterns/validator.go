package tacticalpatterns

// Validator описывает валидатор анемичной модели.
type Validator interface {
	// ValidateModel проверяет валидность модели.
	ValidateModel(model AnemicModel) bool
}

// BusinessEntityValidator описывает валидатор бизнес-сущности.
type BusinessEntityValidator interface {
	Validator
}

// BusinessOperationValidator описывает валидатор бизнес-операции.
type BusinessOperationValidator interface {
	Validator

	// Name возвращает имя бизнес-операции.
	Name() string
}
