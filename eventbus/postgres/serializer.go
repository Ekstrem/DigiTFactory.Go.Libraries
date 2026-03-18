package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
)

// ToEnvelope конвертирует доменное событие в Envelope для записи в outbox.
func ToEnvelope(event *events.DomainEvent) (*Envelope, error) {
	voJSON, err := json.Marshal(event.ChangedValueObjects)
	if err != nil {
		return nil, fmt.Errorf("serialize changed value objects: %w", err)
	}

	return &Envelope{
		BoundedContext:         event.ContextName,
		AggregateID:            event.AggregateID,
		Version:                event.Ver,
		CorrelationToken:       event.Command.CorrToken,
		CommandName:            event.Command.CmdName,
		SubjectName:            event.Command.SubjName,
		CommandVersion:         event.Command.Ver,
		ChangedValueObjectsJSON: string(voJSON),
		Result:                 event.ResultStatus.String(),
		Reason:                 event.ResultReason,
		CreatedAt:              time.Now().UTC(),
	}, nil
}
