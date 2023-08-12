package nikora

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/model"
	"github.com/PuerkitoBio/goquery"
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
		store            model.Store
		discounts        []model.Discount
		discount         model.Discount
		refDiscountsItem map[string]struct{}
	)

	refDiscountsItem = make(map[string]struct{})

	// Page with promotions
	p.Collector.OnHTML("div.cp54.lg_5_5.lp_5_5.md_5_5.sm_12.xs_12", func(h *colly.HTMLElement) {

		ref := fmt.Sprint(p.UrlHost, h.ChildAttr("a", "href"))

		if len(discount.Url) > 0 && discount.Url != ref {
			discounts = append(discounts, discount)
		}

		discount = model.Discount{
			Url:     ref,
			Place:   p.Place,
			Picture: fmt.Sprint(p.UrlHost, h.DOM.Find(".cp3 .cp53 img").AttrOr("src", "")),
			Title:   h.DOM.Find(".cp50 .cp3 h4").Text(),
		}

		h.Request.Visit(discount.Url)

	})

	// Implement pagination
	p.Collector.OnHTML("#pg0", func(h *colly.HTMLElement) {

		nextPage := h.Attr("href")
		// if exists {
		h.Request.Visit(nextPage)
		// }

	})

	// Tag with product
	p.Collector.OnHTML("div.lg_4.lp_6.md_6.sm_6.xs_12.animate_block_when_see.cp100", func(h *colly.HTMLElement) {

		ref := fmt.Sprint(p.UrlHost, h.ChildAttr("a", "href"))
		_, exist := refDiscountsItem[ref]
		if exist {
			return
		}
		refDiscountsItem[ref] = struct{}{}

		regularPrice, err := p.getPrice(h, ".cp17.col_padding .old .lari", ".cp105")
		if err != nil {
			log.Printf("Wrong when get regular price in %s. Err: %v", ref, err)
			return
		}

		finalPrice, err := p.getPrice(h, ".cp17.col_padding .new .lari", ".cp104")
		if err != nil {
			log.Printf("Wrong when get final price in %s. Err: %v", h.ChildAttr("a", "href"), err)
			return
		}

		discountItem := model.DiscountItem{
			Url:          ref,
			Picture:      fmt.Sprint(p.UrlHost, h.DOM.Find(".team.cp7 .cp4 .step1_image_div.cp5 img").AttrOr("src", "")),
			Title:        h.DOM.Find(".team.cp7 .cp4 .cp40 h4").Text(),
			RegularPrice: regularPrice,
			FinalPrice:   finalPrice,
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

func (p *parserDiscount) getPrice(h *colly.HTMLElement, containerSelector string, priceSelector string) (price float64, err error) {

	var (
		intPart   float64
		fractPart float64
	)

	replacer := strings.NewReplacer("\n", "", "\t", "", " ", "")
	containerPrice := h.DOM.Find(containerSelector)

	intPart, err = getPartPrice(replacer, containerPrice, fmt.Sprint(".cp103", priceSelector))
	if err != nil {
		return
	}

	fractPart, err = getPartPrice(replacer, containerPrice, fmt.Sprint(".tetri", priceSelector))
	if err != nil {
		return
	}

	price = intPart + (fractPart / 100)
	price = math.Round(price*100) / 100

	return
}

func getPartPrice(replacer *strings.Replacer, containerPrice *goquery.Selection, selector string) (partPrice float64, err error) {

	attrPartPrice := containerPrice.Find(selector).Text()
	attrPartPrice = replacer.Replace(attrPartPrice)

	partPrice, err = strconv.ParseFloat(attrPartPrice, 64)
	if err != nil {
		return
	}

	return
}
