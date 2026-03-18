package postgres

import (
	"time"

	"github.com/google/uuid"
)

// Envelope — плоский DTO для записи в outbox таблицу.
type Envelope struct {
	BoundedContext         string    `json:"boundedContext"`
	AggregateID            uuid.UUID `json:"aggregateId"`
	Version                int64     `json:"version"`
	CorrelationToken       uuid.UUID `json:"correlationToken"`
	CommandName            string    `json:"commandName"`
	SubjectName            string    `json:"subjectName"`
	CommandVersion         int64     `json:"commandVersion"`
	ChangedValueObjectsJSON string   `json:"changedValueObjectsJson"`
	Result                 string    `json:"result"`
	Reason                 string    `json:"reason"`
	CreatedAt              time.Time `json:"createdAt"`
}
