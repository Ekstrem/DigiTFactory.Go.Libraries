package monads

// Either выполняет разветвление результата.
// Альтернатива тернарному оператору для написания fluent-стиля кода.
func Either[TSource, TResult any](
	source TSource,
	condition func(TSource) bool,
	ifTrue func(TSource) TResult,
	ifFalse func(TSource) TResult,
) TResult {
	if condition(source) {
		return ifTrue(source)
	}
	return ifFalse(source)
}

// PipeTo применяет функцию к значению и возвращает результат.
// Для написания FluentApi-style кода.
func PipeTo[TSource, TResult any](source TSource, fn func(TSource) TResult) TResult {
	return fn(source)
}

// Do выполняет функцию побочного эффекта над объектом и возвращает его.
func Do[T any](obj T, action func(T)) T {
	action(obj)
	return obj
}
