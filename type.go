package paginator

// Direction is sort direction when ordering
type Direction int

const (
	_ Direction = iota
	SortASC
	SortDESC
)

func (enum Direction) String() string {
	switch enum {
	case SortASC:
		return "ASC"
	case SortDESC:
		return "DESC"

	}
	return ""
}

// page contains page number, page size, and total number of items
type Page struct {
	Number int
	Size   int
	Total  int
}

// Order contains sort direction and sort field
type Order struct {
	Column    string
	Direction Direction
}
