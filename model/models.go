package model

import (
	"html/template"

	"golang.org/x/text/language"
)

type Discount struct {
	Url         string
	Place       string
	Picture     string
	Title       string
	Description string
}

type DiscountLanguage struct {
	Url         string
	Place       string
	Picture     string
	Title       string
	Description template.HTML
	Language    language.Tag
}
