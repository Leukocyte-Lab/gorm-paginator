package paginator

import (
	"fmt"

	"gorm.io/gorm"
)

func (pgntr Paginator) Pagging(tx *gorm.DB, dest interface{}) (*Page, error) {
	result := pgntr.GenGormTransaction(tx).Find(&dest)
	if result.Error != nil {
		return nil, fmt.Errorf("Paginator.Pagging : %w", result.Error)
	}

	err := pgntr.CountPageTotal(result)
	if err != nil {
		return nil, fmt.Errorf("Paginator.Pagging : %w", result.Error)
	}

	return &pgntr.Page, nil
}
