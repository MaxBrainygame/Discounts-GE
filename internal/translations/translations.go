package translations

import "github.com/MaxBrainygame/Discounts-GE/model"

type Translator interface {
	TranslateDiscounts(discounts *[]model.Discount) (discountsLanguage *[]model.DiscountLanguage, err error)
}
