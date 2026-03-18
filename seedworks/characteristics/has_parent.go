package characteristics

// HasParent указывает, что объект имеет родительский объект.
type HasParent[TKey comparable] interface {
	// IsHasParent возвращает true, если объект имеет родителя.
	IsHasParent() bool

	// ParentID возвращает идентификатор родителя.
	ParentID() TKey

	// SetParentID устанавливает идентификатор родителя.
	SetParentID(id TKey)
}

// HasNext указывает, что объект имеет следующий объект.
type HasNext[TKey comparable] interface {
	// IsHasNext возвращает true, если объект имеет следующий.
	IsHasNext() bool

	// NextID возвращает идентификатор следующего.
	NextID() TKey

	// SetNextID устанавливает идентификатор следующего.
	SetNextID(id TKey)
}

// HasPrevious указывает, что объект имеет предыдущий объект.
type HasPrevious[TKey comparable] interface {
	// IsHasPrevious возвращает true, если объект имеет предыдущий.
	IsHasPrevious() bool

	// PreviousID возвращает идентификатор предыдущего.
	PreviousID() TKey

	// SetPreviousID устанавливает идентификатор предыдущего.
	SetPreviousID(id TKey)
}
