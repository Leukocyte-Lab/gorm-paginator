package paginator

const (
	MaxPageSize     = 100
	DefaultPageSize = 20
)

// GenPage is helper function to generate page with default page number and page size
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
