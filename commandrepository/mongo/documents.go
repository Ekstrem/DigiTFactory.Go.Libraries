package mongo

import (
	"time"

	"github.com/google/uuid"
)

// DomainEventDocument — документ доменного события в MongoDB.
type DomainEventDocument struct {
	ID                      uuid.UUID `bson:"_id" json:"id"`
	Version                 int64     `bson:"version" json:"version"`
	CorrelationToken        uuid.UUID `bson:"correlation_token" json:"correlationToken"`
	BoundedContext          string    `bson:"bounded_context" json:"boundedContext"`
	CommandName             string    `bson:"command_name" json:"commandName"`
	SubjectName             string    `bson:"subject_name" json:"subjectName"`
	ChangedValueObjectsJSON string    `bson:"changed_value_objects" json:"changedValueObjects"`
	Result                  string    `bson:"result" json:"result"`
	CreatedAt               time.Time `bson:"created_at" json:"createdAt"`
}

// SnapshotDocument — документ snapshot агрегата в MongoDB.
type SnapshotDocument struct {
	ID            uuid.UUID `bson:"_id" json:"id"`
	Version       int64     `bson:"version" json:"version"`
	AggregateJSON string    `bson:"aggregate_json" json:"aggregateJson"`
	CreatedAt     time.Time `bson:"created_at" json:"createdAt"`
}

// AggregateStateDocument — документ текущего состояния агрегата в MongoDB.
type AggregateStateDocument struct {
	ID            uuid.UUID `bson:"_id" json:"id"`
	Version       int64     `bson:"version" json:"version"`
	AggregateJSON string    `bson:"aggregate_json" json:"aggregateJson"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updatedAt"`
}
