package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
)

// Serialize сериализует доменное событие в JSON через Envelope.
func Serialize(event *events.DomainEvent) ([]byte, error) {
	voJSON, err := json.Marshal(event.ChangedValueObjects)
	if err != nil {
		return nil, fmt.Errorf("serialize changed value objects: %w", err)
	}

	envelope := Envelope{
		AggregateID:            event.AggregateID,
		Version:                event.Ver,
		CorrelationToken:       event.Command.CorrToken,
		BoundedContext:         event.ContextName,
		CommandName:            event.Command.CmdName,
		SubjectName:            event.Command.SubjName,
		CommandVersion:         event.Command.Ver,
		ChangedValueObjectsJSON: string(voJSON),
		Result:                 event.ResultStatus.String(),
		Reason:                 event.ResultReason,
		CreatedAt:              time.Now().UTC(),
	}

	return json.Marshal(envelope)
}

// Deserialize десериализует JSON в Envelope.
func Deserialize(data []byte) (*Envelope, error) {
	var envelope Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("deserialize domain event envelope: %w", err)
	}
	return &envelope, nil
}
