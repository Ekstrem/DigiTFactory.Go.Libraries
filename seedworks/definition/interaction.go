package definition

// InteractionOfBoundedContexts определяет тип взаимодействия ограниченных контекстов.
type InteractionOfBoundedContexts int

const (
	// SharedKernel — общее ядро. Часто представляет собой смысловое ядро предметной области (Core Domain),
	// набор естественных подобластей (Generic Subdomains) или то и другое одновременно.
	SharedKernel InteractionOfBoundedContexts = iota

	// Partners — равные отношения взаимовлияния.
	Partners

	// ClientSuppliers — поставщик-потребитель.
	ClientSuppliers

	// Conformist — конформист. Разновидность отношения поставщик-потребитель,
	// которая требует области подстройки. Возникает в том случае,
	// если требования контекста потребителя не учитываются.
	Conformist

	// AntiCorruptionLayer — предохранительный уровень.
	// Разновидность отношения конформист, у которой область подстройки выделена
	// в отдельный слой с полноценной трансляцией для соблюдения правил ограниченного контекста.
	AntiCorruptionLayer

	// OpenHostService — служба с открытым протоколом.
	OpenHostService

	// PublishedLanguage — общедоступный язык.
	PublishedLanguage

	// SeparateWays — отдельное существование.
	SeparateWays
)
