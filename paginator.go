package paginator

import (
	"math"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MinPageNumber = 1
	MinPageSize   = 1
)

type Paginator struct {
	Page   Page
	Order  []Order
	Filter map[string]string
}

func New(page Page, orders []Order, filter map[string]string) *Paginator {
	paginator := Paginator{
		Page:   page,
		Order:  orders,
		Filter: filter,
	}

	paginator.limitMinPageNumber(MinPageNumber)
	paginator.limitMinPageSize(MinPageSize)

	return &paginator
}

// GenGormTransaction: generate GORMv2 sql Transaction (gorm.DB)
func (pgntr Paginator) GenGormTransaction(tx *gorm.DB) *gorm.DB {
	tx = pgntr.offset(tx)
	tx = pgntr.limit(tx)
	tx = pgntr.orderBy(tx)
	tx = pgntr.where(tx)

	return tx
}

// CountPageTotal: setter of Paginator.Page.Total
func (pgntr *Paginator) CountPageTotal(tx *gorm.DB) error {
	var count int64
	// remove offset, limit and order by before count
	delete(tx.Statement.Clauses, "ORDER BY")
	tx.Offset(-1).Limit(-1).Count(&count)
	pgntr.Page.Total = int(math.Ceil(float64(count) / float64(pgntr.Page.Size)))
	// limit PageNumber <= PageTotal
	pgntr.LimitPageTotal()
	return nil
}

func (pgntr Paginator) where(tx *gorm.DB) *gorm.DB {
	return tx.Where(pgntr.Filter)
}

func (pgntr Paginator) offset(tx *gorm.DB) *gorm.DB {
	// concat OFFSET SQL query statement by Paginator.Page.Number/Size
	return tx.Offset((pgntr.Page.Number - 1) * pgntr.Page.Size)
}

func (pgntr Paginator) limit(tx *gorm.DB) *gorm.DB {
	// concat LIMIT SQL query statement by Paginator.Page.Size
	return tx.Limit(pgntr.Page.Size)
}

func (pgntr Paginator) orderBy(tx *gorm.DB) *gorm.DB {
	// concat ORDER SQL query statement by Paginator.order
	for _, order := range pgntr.Order {
		tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: order.Column}, Desc: order.Direction == SortDESC})
	}
	return tx
}

func (pgntr *Paginator) limitMinPageNumber(minPageNumber int) {
	// limit PageNumber >= minPageNumber
	if pgntr.Page.Number < minPageNumber {
		pgntr.Page.Number = minPageNumber
	}
}

func (pgntr *Paginator) limitMinPageSize(minPageSize int) {
	// limit PageSize >= minPageSize
	if pgntr.Page.Size < minPageSize {
		pgntr.Page.Size = minPageSize
	}
}

func (pgntr *Paginator) LimitPageTotal() {
	// set page total default to 1
	if pgntr.Page.Total == 0 {
		pgntr.Page.Total = 1
	}
	// limit PageNumber <= PageTotal
	if pgntr.Page.Number > pgntr.Page.Total {
		pgntr.Page.Number = pgntr.Page.Total
	}
}
