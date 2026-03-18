package characteristics

import "time"

// HasVersion указывает, что объект имеет версионность.
// Тип T определяет тип данных для версионирования.
type HasVersion[T any] interface {
	// Version возвращает версию объекта.
	// Ожидаемое использование — дата создания версии в миллисекундах (Unix timestamp).
	Version() T
}

// HasInt64Version указывает, что объект имеет версионность с типом int64.
type HasInt64Version interface {
	HasVersion[int64]
}

// Versioning предоставляет частичную реализацию версии для передачи в другие параметры.
type Versioning struct {
	stamp int64
}

// NewVersioning создаёт экземпляр Versioning из времени вызова команды.
func NewVersioning(stamp time.Time) Versioning {
	return Versioning{stamp: stamp.UnixMilli()}
}

// Version возвращает дату создания версии в формате Unix-миллисекунд.
func (v Versioning) Version() int64 {
	return v.stamp
}
