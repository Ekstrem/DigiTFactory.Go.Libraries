// Package result содержит типы результатов доменных операций.
package result

// DomainOperationResult определяет результат доменной операции.
type DomainOperationResult int

const (
	// Success — операция выполнена успешно.
	Success DomainOperationResult = iota

	// WithWarnings — операция завершилась, но необходимо отреагировать на несоблюдение инвариантов.
	WithWarnings

	// Exception — операция не выполнена.
	Exception
)

// String возвращает строковое представление результата.
func (r DomainOperationResult) String() string {
	switch r {
	case Success:
		return "Success"
	case WithWarnings:
		return "WithWarnings"
	case Exception:
		return "Exception"
	default:
		return "Unknown"
	}
}

// OperationResult описывает результат выполнения доменной операции.
type OperationResult interface {
	// Result возвращает результат операции.
	Result() DomainOperationResult

	// Reason возвращает причину ошибки.
	Reason() string
}
