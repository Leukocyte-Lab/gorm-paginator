package paginator

const (
	SortASC  = "ASC"
	SortDESC = "DESC"
)

const (
	MinPageNumber = 1
)

const (
	MinPageSize     = 1
	MaxPageSize     = 100
	DefaultPageSize = 20
)

type Page struct {
	Number int
	Size   int
	Total  int
}

type Order struct {
	Column    string
	Direction string
}
