package postgres

import (
	"time"

	"github.com/google/uuid"
)

// DomainEventEntry — запись доменного события в Event Store.
type DomainEventEntry struct {
	// ID — идентификатор агрегата.
	ID uuid.UUID `db:"id" json:"id"`

	// Version — версия события (UTC timestamp в миллисекундах).
	Version int64 `db:"version" json:"version"`

	// CorrelationToken — маркер корреляции для отслеживания цепочки команд.
	CorrelationToken uuid.UUID `db:"correlation_token" json:"correlationToken"`

	// BoundedContext — имя ограниченного контекста.
	BoundedContext string `db:"bounded_context" json:"boundedContext"`

	// CommandName — имя команды, породившей событие.
	CommandName string `db:"command_name" json:"commandName"`

	// SubjectName — имя субъекта, выполнившего команду.
	SubjectName string `db:"subject_name" json:"subjectName"`

	// ChangedValueObjectsJSON — JSON с изменёнными Value Objects.
	ChangedValueObjectsJSON string `db:"changed_value_objects" json:"changedValueObjects"`

	// Result — результат выполнения команды.
	Result string `db:"result" json:"result"`

	// CreatedAt — дата создания события (UTC).
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// SnapshotEntry — snapshot агрегата в Event Store.
type SnapshotEntry struct {
	// ID — идентификатор агрегата.
	ID uuid.UUID `db:"id" json:"id"`

	// Version — версия snapshot.
	Version int64 `db:"version" json:"version"`

	// AggregateJSON — сериализованное состояние агрегата.
	AggregateJSON string `db:"aggregate_json" json:"aggregateJson"`

	// CreatedAt — дата создания snapshot.
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// AggregateStateEntry — текущее состояние агрегата (для стратегии StateOnly).
type AggregateStateEntry struct {
	// ID — идентификатор агрегата.
	ID uuid.UUID `db:"id" json:"id"`

	// Version — текущая версия агрегата.
	Version int64 `db:"version" json:"version"`

	// AggregateJSON — сериализованное состояние агрегата.
	AggregateJSON string `db:"aggregate_json" json:"aggregateJson"`

	// UpdatedAt — дата последнего обновления.
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
