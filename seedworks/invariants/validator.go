package invariants

import (
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/result"
)

// SpecificationValidator — служебная структура валидирования операций.
type SpecificationValidator struct {
	validators map[BusinessOperationSpecification]bool
}

// NewSpecificationValidator создаёт валидатор из набора спецификаций и данных операции.
func NewSpecificationValidator(
	data *result.OperationData,
	specifications ...BusinessOperationSpecification,
) *SpecificationValidator {
	validators := make(map[BusinessOperationSpecification]bool, len(specifications))
	for _, spec := range specifications {
		validators[spec] = spec.IsSatisfiedBy(data)
	}
	return &SpecificationValidator{validators: validators}
}

// Result возвращает true, если все спецификации выполнены.
func (v *SpecificationValidator) Result() bool {
	for _, ok := range v.validators {
		if !ok {
			return false
		}
	}
	return true
}

// GetFailedValidatorsReasons возвращает причины непройденных спецификаций.
func (v *SpecificationValidator) GetFailedValidatorsReasons() map[string]result.DomainOperationResult {
	failed := make(map[string]result.DomainOperationResult)
	for spec, ok := range v.validators {
		if !ok {
			failed[spec.Reason()] = spec.DomainResult()
		}
	}
	return failed
}
