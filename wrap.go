package paginator

import (
	"fmt"

	"gorm.io/gorm"
)

// Paging : A easy use wrapper for doing paging
func (pgntr Paginator) Paging(tx *gorm.DB, dest interface{}) (*Page, error) {
	result := pgntr.GenGormTransaction(tx).Find(&dest)
	if result.Error != nil {
		return nil, fmt.Errorf("paginator.paging : %w", result.Error)
	}

	err := pgntr.CountPageTotal(result)
	if err != nil {
		return nil, fmt.Errorf("paginator.paging : %w", result.Error)
	}

	return &pgntr.Page, nil
}
