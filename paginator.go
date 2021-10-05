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

// Paginator is basic struct contains pagination information
type Paginator struct {
	Page   Page
	Order  []Order
	Filter map[string]string
}

// New is helper function for create Paginator instance
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

// GenGormTransaction generate GORMv2 sql Transaction (gorm.DB) for paging
func (pgntr Paginator) GenGormTransaction(tx *gorm.DB) *gorm.DB {
	tx = pgntr.GenPgntrStmt(tx)
	tx = pgntr.where(tx)

	return tx
}

// CountPageTotal is setter of Paginator.Page.Total by counting total page number
func (pgntr *Paginator) CountPageTotal(tx *gorm.DB) error {
	// remove offset, limit and order by before count
	tx = pgntr.DelPgntrStmt(tx)

	var recordCount int64
	tx.Count(&recordCount)

	// set PageTotal by counting total page number
	pgntr.Page.Total = pgntr.countPageTotal(recordCount)
	pgntr.LimitPageTotal()

	return nil
}

// GenPgntrStmt add statements for paging
func (pgntr Paginator) GenPgntrStmt(tx *gorm.DB) *gorm.DB {
	tx = pgntr.offset(tx)
	tx = pgntr.limit(tx)
	tx = pgntr.orderBy(tx)

	return tx
}

// DelPgntrStmt remove statements for paging
func (pgntr Paginator) DelPgntrStmt(tx *gorm.DB) *gorm.DB {
	tx = pgntr.rmOrderBy(tx)
	tx = pgntr.rmLimit(tx)
	tx = pgntr.rmOffset(tx)

	return tx
}

func (pgntr Paginator) where(tx *gorm.DB) *gorm.DB {
	return tx.Where(pgntr.Filter)
}

func (pgntr Paginator) offset(tx *gorm.DB) *gorm.DB {
	// concat OFFSET SQL query statement by Paginator.Page.Number/Size
	return tx.Offset((pgntr.Page.Number - 1) * pgntr.Page.Size)
}

func (pgntr Paginator) rmOffset(tx *gorm.DB) *gorm.DB {
	return tx.Offset(-1)
}

func (pgntr Paginator) limit(tx *gorm.DB) *gorm.DB {
	// concat LIMIT SQL query statement by Paginator.Page.Size
	return tx.Limit(pgntr.Page.Size)
}

func (pgntr Paginator) rmLimit(tx *gorm.DB) *gorm.DB {
	return tx.Limit(-1)
}

func (pgntr Paginator) orderBy(tx *gorm.DB) *gorm.DB {
	// concat ORDER SQL query statement by Paginator.order
	for _, order := range pgntr.Order {
		tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: order.Column}, Desc: order.Direction == SortDESC})
	}
	return tx
}

func (pgntr Paginator) rmOrderBy(tx *gorm.DB) *gorm.DB {
	delete(tx.Statement.Clauses, "ORDER BY")
	return tx
}

func (pgntr Paginator) countPageTotal(recordCount int64) (pageTotal int) {
	// page total start with 1
	if recordCount == 0 || pgntr.Page.Size == 0 {
		return 1
	}

	return int(math.Ceil(float64(recordCount) / float64(pgntr.Page.Size)))
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
	pgntr.limitPageTotal()
	pgntr.limitPageNumber()
}

func (pgntr *Paginator) limitPageTotal() {
	// set page total default to 1
	if pgntr.Page.Total == 0 {
		pgntr.Page.Total = 1
	}
}

func (pgntr *Paginator) limitPageNumber() {
	// limit PageNumber <= PageTotal
	if pgntr.Page.Number > pgntr.Page.Total {
		pgntr.Page.Number = pgntr.Page.Total
	}
}
