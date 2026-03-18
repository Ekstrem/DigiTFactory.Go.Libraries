package invariants

import (
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/definition"
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/result"
)

// ValidateCommand выполняет валидацию бизнес-операции через набор спецификаций.
// Сначала проверяются утверждения (assertions), затем предупреждения (validators).
// Возвращает соответствующий результат: Exception, WithWarnings или Success.
func ValidateCommand(
	data *result.OperationData,
	bc definition.BoundedContextDescription,
	specifications ...BusinessOperationSpecification,
) *result.AggregateResult {
	// Разделяем assertions и validators
	var assertions []BusinessOperationSpecification
	var validators []BusinessOperationSpecification

	for _, spec := range specifications {
		if _, ok := spec.(BusinessOperationAssertion); ok {
			assertions = append(assertions, spec)
		} else {
			validators = append(validators, spec)
		}
	}

	// Проверяем assertions (ошибки — операция невозможна)
	if len(assertions) > 0 {
		v := NewSpecificationValidator(data, assertions...)
		if !v.Result() {
			reasons := failedReasons(v, result.Exception)
			return result.NewAggregateResultException(data, bc, reasons)
		}
	}

	// Проверяем validators (предупреждения — операция возможна, но с замечаниями)
	if len(validators) > 0 {
		v := NewSpecificationValidator(data, validators...)
		if !v.Result() {
			reasons := failedReasons(v, result.WithWarnings)
			return result.NewAggregateResultWithWarnings(data, bc, reasons)
		}
	}

	return result.NewAggregateResultSuccess(data, bc)
}

func failedReasons(v *SpecificationValidator, level result.DomainOperationResult) []string {
	var reasons []string
	for reason, domainResult := range v.GetFailedValidatorsReasons() {
		if domainResult == level {
			reasons = append(reasons, reason)
		}
	}
	return reasons
}
