// Package reactive содержит реактивные паттерны: управление подписками.
package reactive

// Unsubscriber управляет отпиской наблюдателя из списка наблюдателей.
// Реализует паттерн IDisposable из C# через метод Close.
type Unsubscriber[T any] struct {
	observers *[]T
	observer  T
	index     int
}

// NewUnsubscriber создаёт объект отписки и добавляет наблюдателя в список.
func NewUnsubscriber[T comparable](observers *[]T, observer T) *Unsubscriber[T] {
	// Проверяем, есть ли уже наблюдатель
	for i, obs := range *observers {
		if obs == observer {
			return &Unsubscriber[T]{observers: observers, observer: observer, index: i}
		}
	}

	*observers = append(*observers, observer)
	return &Unsubscriber[T]{
		observers: observers,
		observer:  observer,
		index:     len(*observers) - 1,
	}
}

// Close отписывает наблюдателя из списка.
func (u *Unsubscriber[T]) Close() error {
	if u.observers == nil {
		return nil
	}
	obs := *u.observers
	if u.index >= 0 && u.index < len(obs) {
		*u.observers = append(obs[:u.index], obs[u.index+1:]...)
	}
	return nil
}
