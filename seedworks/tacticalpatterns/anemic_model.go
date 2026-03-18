package tacticalpatterns

import (
	"github.com/google/uuid"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/characteristics"
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/result"
)

// AnemicModel описывает анемичную модель ограниченного контекста для фабрики создания агрегата.
// Анемичная модель должна содержать объект-значения.
type AnemicModel interface {
	characteristics.HasGuidKey
	characteristics.HasInt64Version
	characteristics.CommandSubject
	characteristics.HasGuidCorrelationToken
	characteristics.HasValueObjects

	// Также реализует AnemicModelView для использования в result
	result.AnemicModelView
}

// BaseAnemicModel — базовая реализация анемичной модели.
type BaseAnemicModel struct {
	AggregateID      uuid.UUID      `json:"id"`
	AggregateVersion int64          `json:"version"`
	CorrToken        uuid.UUID      `json:"correlationToken"`
	ValueObjects     map[string]any `json:"invariants"`
}

// NewBaseAnemicModel создаёт анемичную модель из комплексного ключа и набора объект-значений.
func NewBaseAnemicModel(key characteristics.ComplexKey, invariants map[string]any) BaseAnemicModel {
	if invariants == nil {
		invariants = make(map[string]any)
	}
	return BaseAnemicModel{
		AggregateID:      key.ID(),
		AggregateVersion: key.Version(),
		CorrToken:        key.CorrelationToken(),
		ValueObjects:     invariants,
	}
}

// ID возвращает идентификатор сущности.
func (m BaseAnemicModel) ID() uuid.UUID { return m.AggregateID }

// Version возвращает версию агрегата.
func (m BaseAnemicModel) Version() int64 { return m.AggregateVersion }

// CorrelationToken возвращает маркер корреляции.
func (m BaseAnemicModel) CorrelationToken() uuid.UUID { return m.CorrToken }

// CommandName возвращает имя команды (по умолчанию пустая строка).
func (m BaseAnemicModel) CommandName() string { return "" }

// SubjectName возвращает имя субъекта бизнес-операции (по умолчанию пустая строка).
func (m BaseAnemicModel) SubjectName() string { return "" }

// Invariants возвращает словарь объект-значений (инвариантов) модели.
func (m BaseAnemicModel) Invariants() map[string]any { return m.ValueObjects }

// GetValueObjects возвращает объект-значения модели.
func (m BaseAnemicModel) GetValueObjects() map[string]any { return m.ValueObjects }
