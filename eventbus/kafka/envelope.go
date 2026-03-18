package kafka

import (
	"time"

	"github.com/google/uuid"
)

// Envelope — плоский DTO для сериализации доменного события в JSON.
type Envelope struct {
	AggregateID            uuid.UUID `json:"aggregateId"`
	Version                int64     `json:"version"`
	CorrelationToken       uuid.UUID `json:"correlationToken"`
	BoundedContext         string    `json:"boundedContext"`
	CommandName            string    `json:"commandName"`
	SubjectName            string    `json:"subjectName"`
	CommandVersion         int64     `json:"commandVersion"`
	ChangedValueObjectsJSON string   `json:"changedValueObjectsJson"`
	Result                 string    `json:"result"`
	Reason                 string    `json:"reason"`
	CreatedAt              time.Time `json:"createdAt"`
}
