package parsers

import (
	"github.com/MaxBrainygame/Discounts-GE/model"
)

type ParseDiscounter interface {
	ParseDiscounts() (discounts *[]model.Discount, err error)
}
