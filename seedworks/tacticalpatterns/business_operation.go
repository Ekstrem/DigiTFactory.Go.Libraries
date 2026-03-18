package tacticalpatterns

import (
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/result"
)

// AggregateBusinessOperation описывает бизнес-операцию над агрегатом.
type AggregateBusinessOperation interface {
	// Handle выполняет бизнес-операцию.
	Handle(
		model AnemicModel,
		command *events.CommandToAggregateData,
		scope BoundedContextScope,
	) *result.AggregateResult
}
