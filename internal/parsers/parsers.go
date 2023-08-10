package parsers

import (
	"github.com/MaxBrainygame/Discounts-GE/model"
)

type ParseDiscounter interface {
	ParseDiscounts(categoryStores map[string]*model.CategoryStores) (discounts *model.Store, err error)
}
