package parsers

import (
	"github.com/MaxBrainygame/Discounts-GE/model"
)

type ParseDiscounter interface {
	ParseDiscounts() (discounts *model.Store, err error)
}
