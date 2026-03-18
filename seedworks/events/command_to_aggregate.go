// Package events содержит типы доменных событий и интерфейсы шины событий.
package events

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CommandToAggregateData содержит сведения о команде к агрегату.
type CommandToAggregateData struct {
	// CorrToken — маркер корреляции.
	CorrToken uuid.UUID `json:"correlationToken"`

	// Ver — версия (Unix-миллисекунды).
	Ver int64 `json:"version"`

	// CmdName — имя метода агрегата, который вызывает команда.
	CmdName string `json:"commandName"`

	// SubjName — имя субъекта бизнес-операции.
	SubjName string `json:"subjectName"`
}

// CorrelationToken возвращает маркер корреляции.
func (c *CommandToAggregateData) CorrelationToken() uuid.UUID { return c.CorrToken }

// Version возвращает версию (Unix-миллисекунды).
func (c *CommandToAggregateData) Version() int64 { return c.Ver }

// CommandName возвращает имя метода агрегата, который вызывает команда.
func (c *CommandToAggregateData) CommandName() string { return c.CmdName }

// SubjectName возвращает имя субъекта бизнес-операции.
func (c *CommandToAggregateData) SubjectName() string { return c.SubjName }

// String возвращает строковое представление команды.
func (c *CommandToAggregateData) String() string {
	return fmt.Sprintf("%s %s", c.CmdName, c.SubjName)
}

// NewCommandToAggregate создаёт структуру CommandToAggregateData.
func NewCommandToAggregate(
	correlationToken uuid.UUID,
	commandName string,
	subjectName string,
	version int64,
) *CommandToAggregateData {
	return &CommandToAggregateData{
		CorrToken: correlationToken,
		Ver:       version,
		CmdName:   commandName,
		SubjName:  subjectName,
	}
}

// NewCommandToAggregateNow создаёт структуру CommandToAggregateData
// с версией по умолчанию (текущее время в Unix-миллисекундах).
func NewCommandToAggregateNow(
	correlationToken uuid.UUID,
	commandName string,
	subjectName string,
) *CommandToAggregateData {
	return &CommandToAggregateData{
		CorrToken: correlationToken,
		Ver:       time.Now().UTC().UnixMilli(),
		CmdName:   commandName,
		SubjName:  subjectName,
	}
}
