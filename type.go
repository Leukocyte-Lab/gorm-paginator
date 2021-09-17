package paginator

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

type Page struct {
	Number int
	Size   int
	Total  int
}

type Order struct {
	Column    string
	Direction Direction
}
