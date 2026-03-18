package repository

import "context"

// SqlProgrammability описывает вызов хранимых процедур и функций SQL.
type SqlProgrammability interface {
	// ExecScalarFunc выполняет скалярную функцию.
	ExecScalarFunc(ctx context.Context, sqlFuncName string, params ...any) (any, error)

	// ExecTableFunc выполняет табличную функцию.
	ExecTableFunc(ctx context.Context, sqlFuncName string, result any, params ...any) error

	// ExecStoredProc выполняет хранимую процедуру.
	ExecStoredProc(ctx context.Context, sqlProcName string, params ...any) error

	// ExecStoredProcList выполняет хранимую процедуру и возвращает список результатов.
	ExecStoredProcList(ctx context.Context, sqlProcName string, result any, params ...any) error

	// ExecStoredProcSingle выполняет хранимую процедуру и возвращает единичный результат.
	ExecStoredProcSingle(ctx context.Context, sqlProcName string, result any, params ...any) error
}
