package characteristics

// Paging описывает интерфейс постраничной навигации.
type Paging interface {
	// Page возвращает номер страницы.
	Page() int

	// PageSize возвращает количество элементов на странице.
	PageSize() int
}

// PagingParams — реализация интерфейса Paging.
type PagingParams struct {
	page     int
	pageSize int
}

// NewPaging создаёт экземпляр PagingParams.
func NewPaging(page, pageSize int) PagingParams {
	return PagingParams{page: page, pageSize: pageSize}
}

// Page возвращает номер страницы.
func (p PagingParams) Page() int { return p.page }

// PageSize возвращает количество элементов на странице.
func (p PagingParams) PageSize() int { return p.pageSize }
