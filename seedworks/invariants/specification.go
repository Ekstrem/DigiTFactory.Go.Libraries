// Package invariants содержит спецификации и валидаторы бизнес-операций.
package invariants

import (
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/result"
)

// BusinessOperationSpecification описывает спецификацию бизнес-операции.
// Проверяет, что результат операции соблюдает инварианты и может быть применён.
type BusinessOperationSpecification interface {
	// IsSatisfiedBy проверяет соответствие спецификации.
	IsSatisfiedBy(data *result.OperationData) bool

	// Reason возвращает причину, по которой объект не проходит проверку спецификацией.
	Reason() string

	// DomainResult возвращает результат бизнес-операции в случае непройденной спецификации.
	DomainResult() result.DomainOperationResult
}

// BusinessOperationAssertion описывает спецификацию-утверждение.
// Проверяет, что результат операции непротиворечив и результат нельзя применить
// в случае нарушения утверждения.
type BusinessOperationAssertion interface {
	BusinessOperationSpecification
}

// BusinessOperationValidatorSpec описывает спецификацию проверки на возможность
// проведения бизнес-операции.
type BusinessOperationValidatorSpec interface {
	BusinessOperationSpecification
}
