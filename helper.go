package paginator

const (
	MaxPageSize     = 100
	DefaultPageSize = 20
)

func GenPage(pageNo int, pageSize int) Page {
	if pageNo < MinPageNumber {
		pageNo = MinPageNumber
	}
	switch {
	case pageSize > MaxPageSize:
		pageSize = MaxPageSize
	case pageSize < MinPageSize:
		pageSize = DefaultPageSize
	}

	return Page{
		Number: pageNo,
		Size:   pageSize,
	}
}
