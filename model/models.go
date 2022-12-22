package model

import (
	"golang.org/x/text/language"
)

type Discount struct {
	Url          string
	Place        string
	Picture      string
	Title        string
	Description  string
	RegularPrice float64
	FinalPrice   float64
}

type DiscountLanguage struct {
	Url          string
	Place        string
	Picture      string
	Title        string
	RegularPrice float64
	FinalPrice   float64
	Description  string
	Language     language.Tag
}
