package repository

// SingleReadMode определяет режим чтения единичной записи.
type SingleReadMode int

const (
	// Single — единственная запись (вернёт ошибку, если записей больше одной).
	Single SingleReadMode = iota

	// First — первая запись.
	First

	// Last — последняя запись.
	Last
)
