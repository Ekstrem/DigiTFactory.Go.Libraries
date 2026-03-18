package characteristics

// HasValueObjects указывает, что объект содержит объект-значения (инварианты).
type HasValueObjects interface {
	// Invariants возвращает словарь объект-значений (инвариантов).
	Invariants() map[string]any

	// GetValueObjects возвращает объект-значения модели.
	GetValueObjects() map[string]any
}
