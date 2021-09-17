package paginator

const (
	SortASC  = "ASC"
	SortDESC = "DESC"
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
