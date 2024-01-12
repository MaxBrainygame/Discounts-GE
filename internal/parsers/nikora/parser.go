package nikora

import (
	"fmt"
	"strings"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/model"
	"github.com/gocolly/colly"
)

const (
	urlHost     = "http://nikorasupermarket.ge"
	discountUrl = "/ge/მიმდინარე-აქციები"
	place       = "Nikora"
	logoUrl     = "https://nikorasupermarket.ge/images/logo.svg"
	typeStore   = "Grocery"
)

type parserDiscount struct {
	Collector *colly.Collector
	Place     string
	Url       string
	UrlHost   string
	LogoUrl   string
	TypeStore string
}

func NewParser() parsers.ParseDiscounter {

	return &parserDiscount{
		Collector: colly.NewCollector(),
		Place:     place,
		Url:       fmt.Sprint(urlHost, discountUrl),
		UrlHost:   urlHost,
		LogoUrl:   logoUrl,
		TypeStore: typeStore,
	}

}

func (p *parserDiscount) ParseDiscounts(categoryStores map[string]*model.CategoryStores) (*model.Store, error) {

	var (
		store     model.Store
		discounts []model.Discount
		discount  model.Discount
	)

	// Page with promotions
	p.Collector.OnHTML("div.cp54.lg_5_5.lp_5_5.md_5_5.sm_12.xs_12", func(h *colly.HTMLElement) {

		ref := fmt.Sprint(p.UrlHost, h.DOM.Find(".cp50 .pdf").AttrOr("href", ""))

		if len(discount.Url) > 0 && discount.Url != ref {
			discounts = append(discounts, discount)
		}

		discount = model.Discount{
			Url:     ref,
			Place:   p.Place,
			Picture: fmt.Sprint(p.UrlHost, h.DOM.Find(".cp3 .cp53 img").AttrOr("src", "")),
			Title:   h.DOM.Find(".cp50 .cp3 h4").Text(),
		}

		refItem := ""
		indexStr := strings.Index(discount.Url, "=")
		if indexStr > -1 {
			refItem = fmt.Sprint(p.UrlHost, discount.Url[indexStr+1:len(discount.Url)])
		}

		discountItem := model.DiscountItem{
			Url:          refItem,
			Picture:      "",
			Title:        discount.Title,
			RegularPrice: 0,
			FinalPrice:   0,
		}

		discount.Goods = append(discount.Goods, discountItem)

	})

	p.Collector.Visit(p.Url)

	// Add last element
	discounts = append(discounts, discount)

	store = model.Store{
		Name:      p.Place,
		Logo:      p.LogoUrl,
		Host:      p.UrlHost,
		Category:  *categoryStores[p.TypeStore],
		Discounts: discounts,
	}

	return &store, nil

}
