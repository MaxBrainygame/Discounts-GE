package model

import (
	"golang.org/x/text/language"
)

type Discount struct {
	Url         string
	Place       string
	Picture     string
	Title       string
	Description string
	Goods       []DiscountItem
}

type DiscountLanguage struct {
	Url         string
	Place       string
	Picture     string
	Title       string
	Description string
	Language    language.Tag
	Goods       []DiscountItem
}

type DiscountItem struct {
	Url          string
	Picture      string
	Title        string
	RegularPrice float64
	FinalPrice   float64
}
