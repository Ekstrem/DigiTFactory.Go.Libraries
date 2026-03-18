package result

import "github.com/google/uuid"

// AnemicModelView предоставляет минимальный набор данных анемичной модели,
// необходимый для построения результата операции.
// Позволяет избежать циклической зависимости с пакетом tacticalpatterns.
type AnemicModelView interface {
	// ID возвращает идентификатор агрегата.
	ID() uuid.UUID

	// Version возвращает версию агрегата.
	Version() int64

	// CorrelationToken возвращает маркер корреляции.
	CorrelationToken() uuid.UUID

	// CommandName возвращает имя команды.
	CommandName() string

	// SubjectName возвращает имя субъекта бизнес-операции.
	SubjectName() string

	// GetValueObjects возвращает словарь объект-значений.
	GetValueObjects() map[string]any
}

// OperationData содержит данные бизнес-операции: агрегат (до) и модель (после).
type OperationData struct {
	// Aggregate — состояние агрегата до выполнения операции.
	Aggregate AnemicModelView

	// Model — состояние модели после выполнения операции.
	Model AnemicModelView
}

// NewOperationData создаёт экземпляр OperationData.
func NewOperationData(aggregate, model AnemicModelView) *OperationData {
	return &OperationData{
		Aggregate: aggregate,
		Model:     model,
	}
}

// GetChangedValueObjects возвращает словарь изменившихся объект-значений.
// Сравнивает объект-значения агрегата (до) и модели (после) и возвращает только различающиеся.
func (d *OperationData) GetChangedValueObjects() map[string]any {
	aggregateVOs := d.Aggregate.GetValueObjects()
	modelVOs := d.Model.GetValueObjects()

	// Собираем все ключи
	allKeys := make(map[string]struct{})
	for k := range aggregateVOs {
		allKeys[k] = struct{}{}
	}
	for k := range modelVOs {
		allKeys[k] = struct{}{}
	}

	changed := make(map[string]any)
	for key := range allKeys {
		aggVal, aggOk := aggregateVOs[key]
		modVal, modOk := modelVOs[key]

		// Если значения различаются и новое значение существует
		if modOk && (!aggOk || aggVal != modVal) {
			changed[key] = modVal
		}
	}

	return changed
}
