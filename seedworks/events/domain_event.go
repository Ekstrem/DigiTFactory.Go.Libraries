package events

import (
	"github.com/google/uuid"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/definition"
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/result"
)

// DomainEvent — доменное событие предметной области.
type DomainEvent struct {
	// AggregateID — идентификатор экземпляра агрегата.
	AggregateID uuid.UUID `json:"id"`

	// Ver — версия агрегата, над которой выполнялась команда.
	Ver int64 `json:"version"`

	// Command — имя бизнес-операции (доменное событие).
	Command *CommandToAggregateData `json:"command"`

	// ChangedValueObjects — словарь изменившихся объект-значений.
	ChangedValueObjects map[string]any `json:"changedValueObjects"`

	// BoundedContextName — имя ограниченного контекста, в котором произошло событие.
	BoundedContextName string `json:"boundedContext"`

	// ContextName — имя ограниченного контекста.
	ContextName string `json:"contextName"`

	// MsVersion — версия микросервиса.
	MsVersion int `json:"microserviceVersion"`

	// ResultStatus — результат доменной операции.
	ResultStatus result.DomainOperationResult `json:"result"`

	// ResultReason — причина ошибки при выполнении доменной операции.
	ResultReason string `json:"reason"`
}

// ID возвращает идентификатор агрегата.
func (e *DomainEvent) ID() uuid.UUID { return e.AggregateID }

// Version возвращает версию агрегата.
func (e *DomainEvent) Version() int64 { return e.Ver }

// CorrelationToken возвращает маркер корреляции команды.
func (e *DomainEvent) CorrelationToken() uuid.UUID { return e.Command.CorrToken }

// Result возвращает результат доменной операции.
func (e *DomainEvent) Result() result.DomainOperationResult { return e.ResultStatus }

// Reason возвращает причину ошибки при выполнении доменной операции.
func (e *DomainEvent) Reason() string { return e.ResultReason }

// NewDomainEvent создаёт доменное событие.
func NewDomainEvent(
	id uuid.UUID,
	version int64,
	command *CommandToAggregateData,
	changedValueObjects map[string]any,
	description definition.BoundedContextDescription,
	resultStatus result.DomainOperationResult,
	reason string,
) *DomainEvent {
	return &DomainEvent{
		AggregateID:         id,
		Ver:                 version,
		Command:             command,
		ChangedValueObjects: changedValueObjects,
		BoundedContextName:  description.ContextName,
		ContextName:         description.ContextName,
		MsVersion:           description.MicroserviceVersion,
		ResultStatus:        resultStatus,
		ResultReason:        reason,
	}
}

// NewDomainEventFromResult создаёт доменное событие из результата бизнес-операции.
func NewDomainEventFromResult(r *result.AggregateResult) *DomainEvent {
	reason := ""
	if len(r.Reasons) > 0 {
		reason = r.Reasons[0]
	}

	cmd := NewCommandToAggregate(
		r.OperationData.Model.CorrelationToken(),
		r.OperationData.Model.CommandName(),
		r.OperationData.Model.SubjectName(),
		r.OperationData.Model.Version(),
	)

	return NewDomainEvent(
		r.OperationData.Aggregate.ID(),
		r.OperationData.Aggregate.Version(),
		cmd,
		r.OperationData.GetChangedValueObjects(),
		r.BoundedContext,
		r.ResultStatus,
		reason,
	)
}
