package repository

import "context"

// UnitOfWork описывает единицу работы (Unit of Work).
type UnitOfWork[T any] interface {
	// GetCommandRepository возвращает репозиторий команд.
	GetCommandRepository() CommandRepository[T]

	// GetQueryRepository возвращает репозиторий запросов.
	GetQueryRepository() QueryRepository[T]

	// Save сохраняет изменения синхронно.
	Save() (int, error)

	// SaveAsync сохраняет изменения асинхронно.
	SaveAsync(ctx context.Context) (int, error)

	// Programmability возвращает доступ к хранимым процедурам.
	Programmability() SqlProgrammability

	// TransactionManager возвращает управление транзакциями.
	TransactionManager() Transaction
}
