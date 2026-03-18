// Package commandrepository содержит общие типы для Event Store.
package commandrepository

// EventStoreStrategy определяет стратегию хранения событий в Event Store.
type EventStoreStrategy int

const (
	// FullEventSourcing — все события сохраняются, агрегат восстанавливается из полного стрима.
	FullEventSourcing EventStoreStrategy = iota

	// SnapshotAfterN — каждые N событий сохраняется snapshot агрегата.
	// При чтении: snapshot + события после него.
	SnapshotAfterN

	// StateOnly — каждый раз сохраняется агрегат целиком (без истории событий).
	StateOnly
)
