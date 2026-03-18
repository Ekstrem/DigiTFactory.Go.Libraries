// Package monads содержит функциональные конструкции (монады).
package monads

import "fmt"

// Result представляет результат обработки, который может быть успешным или неудачным.
type Result[TSuccess, TFailure any] struct {
	success   TSuccess
	failure   TFailure
	isSuccess bool
}

// Match выполняет разветвление по результату операции.
func (r Result[TSuccess, TFailure]) Match(
	success func(TSuccess) any,
	failure func(TFailure) any,
) any {
	if r.isSuccess {
		return success(r.success)
	}
	return failure(r.failure)
}

// MatchTyped выполняет типизированное разветвление по результату операции.
func MatchTyped[TSuccess, TFailure, TResult any](
	r Result[TSuccess, TFailure],
	success func(TSuccess) TResult,
	failure func(TFailure) TResult,
) TResult {
	if r.isSuccess {
		return success(r.success)
	}
	return failure(r.failure)
}

// IsSuccess возвращает true, если операция завершилась успешно.
func (r Result[TSuccess, TFailure]) IsSuccess() bool {
	return r.isSuccess
}

// String возвращает строковое представление результата.
func (r Result[TSuccess, TFailure]) String() string {
	return fmt.Sprintf("is Success: %t", r.isSuccess)
}

// Success возвращает успешное значение.
// Паникует, если результат — неудача.
func (r Result[TSuccess, TFailure]) Success() TSuccess {
	if !r.isSuccess {
		panic("attempt to get success value from failure result")
	}
	return r.success
}

// Failure возвращает значение ошибки.
// Паникует, если результат — успех.
func (r Result[TSuccess, TFailure]) Failure() TFailure {
	if r.isSuccess {
		panic("attempt to get failure value from success result")
	}
	return r.failure
}

// NewSuccess создаёт успешный результат.
func NewSuccess[TSuccess, TFailure any](success TSuccess) Result[TSuccess, TFailure] {
	return Result[TSuccess, TFailure]{success: success, isSuccess: true}
}

// NewFailure создаёт неудачный результат.
func NewFailure[TSuccess, TFailure any](failure TFailure) Result[TSuccess, TFailure] {
	return Result[TSuccess, TFailure]{failure: failure, isSuccess: false}
}
