package repository

import "context"

// RebuildService описывает сервис перестроения Read-модели (проекций) из Event Store.
// Читает все события в хронологическом порядке и прогоняет через Projection Handler.
// Идемпотентен: повторный вызов безопасен (события с Version <= текущей пропускаются).
type RebuildService interface {
	// RebuildAsync перестраивает все проекции из Event Store.
	RebuildAsync(ctx context.Context) error
}
