package result

import (
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/definition"
)

// AggregateResult описывает результат выполнения бизнес-операции в агрегате.
type AggregateResult struct {
	// OperationData содержит данные бизнес-операции (агрегат до и после).
	OperationData *OperationData

	// BoundedContext — описание ограниченного контекста.
	BoundedContext definition.BoundedContextDescription

	// ResultStatus — результат операции.
	ResultStatus DomainOperationResult

	// Reasons — причины в случае неуспеха выполнения операции.
	Reasons []string
}

// Result возвращает результат операции.
func (r *AggregateResult) Result() DomainOperationResult {
	return r.ResultStatus
}

// Reason возвращает первую причину ошибки или пустую строку.
func (r *AggregateResult) Reason() string {
	if len(r.Reasons) > 0 {
		return r.Reasons[0]
	}
	return ""
}

// NewAggregateResultSuccess создаёт успешный результат операции.
func NewAggregateResultSuccess(data *OperationData, bc definition.BoundedContextDescription) *AggregateResult {
	return &AggregateResult{
		OperationData:  data,
		BoundedContext: bc,
		ResultStatus:   Success,
		Reasons:        nil,
	}
}

// NewAggregateResultException создаёт результат операции с ошибкой.
func NewAggregateResultException(data *OperationData, bc definition.BoundedContextDescription, reasons []string) *AggregateResult {
	return &AggregateResult{
		OperationData:  data,
		BoundedContext: bc,
		ResultStatus:   Exception,
		Reasons:        reasons,
	}
}

// NewAggregateResultWithWarnings создаёт результат операции с предупреждениями.
func NewAggregateResultWithWarnings(data *OperationData, bc definition.BoundedContextDescription, reasons []string) *AggregateResult {
	return &AggregateResult{
		OperationData:  data,
		BoundedContext: bc,
		ResultStatus:   WithWarnings,
		Reasons:        reasons,
	}
}
