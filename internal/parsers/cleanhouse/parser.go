package cleanhouse

import (
	"fmt"
	"log"
	"strconv"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/model"
	"github.com/gocolly/colly"
)

const (
	urlresource = "https://ch.ge"
	discountUrl = "/promotions-list.html"
	place       = "CleanHouse"
)

type parserDiscount struct {
	Collector *colly.Collector
	Place     string
}

func NewParser() parsers.ParseDiscounter {

	return &parserDiscount{
		Collector: colly.NewCollector(),
		Place:     place,
	}
}

func (p *parserDiscount) ParseDiscounts() (*[]model.Discount, error) {

	var (
		discounts []model.Discount
		discount  model.Discount
	)

	// Page with promotions
	p.Collector.OnHTML("div.ab__dotd_promotions-item", func(h *colly.HTMLElement) {

		ref := h.ChildAttr("a", "href")

		if len(discount.Url) > 0 && discount.Url != ref {
			discounts = append(discounts, discount)
		}

		discount = model.Discount{
			Url:     ref,
			Place:   p.Place,
			Picture: h.DOM.Find(".ab__dotd_promotions-item_image a img").AttrOr("data-src", ""),
			Title:   h.DOM.Find("div.ab__dotd_promotions-item_title").Text(),
		}

		h.Request.Visit(ref)

	})

	// Implement pagination
	p.Collector.OnHTML("div.ty-pagination", func(h *colly.HTMLElement) {

		nextPage, exists := h.DOM.Find(".ty-pagination__item.ty-pagination__btn.ty-pagination__next.cm-history.cm-ajax.ty-pagination__right-arrow").Attr("href")
		if exists {
			h.Request.Visit(nextPage)
		}

	})

	// Tag with product
	p.Collector.OnHTML("div.ut2-gl__body", func(h *colly.HTMLElement) {

		regularPrice, err := getPrice(h, ".ut2-gl__price div span .ty-list-price.ty-nowrap .ty-strike span")
		if err != nil {
			log.Printf("Wrong when get regular price in %s. Err: %v", h.ChildAttr("a", "href"), err)
			return
		}
		finalPrice, err := getPrice(h, ".ut2-gl__price div span h2.ty-price span")
		if err != nil {
			log.Printf("Wrong when get final price in %s. Err: %v", h.ChildAttr("a", "href"), err)
			return
		}

		discountItem := model.DiscountItem{
			Url: h.ChildAttr("a", "href"),
			//Place:        place,
			Picture:      getPicture(h),
			Title:        h.DOM.Find("h4.ut2-gl__name").Text(),
			RegularPrice: regularPrice,
			FinalPrice:   finalPrice,
		}

		discount.Goods = append(discount.Goods, discountItem)
	})

	p.Collector.Visit(fmt.Sprintf("%v%v", urlresource, discountUrl))

	return &discounts, nil
}

func getPrice(h *colly.HTMLElement, selector string) (price float64, err error) {

	var attrPrice string

	containerPrice := h.DOM.Find(selector)
	if len(containerPrice.Nodes) > 0 {
		attrPrice = containerPrice.Nodes[0].FirstChild.Data
	}

	price, err = strconv.ParseFloat(attrPrice, 64)

	return
}

func getPicture(h *colly.HTMLElement) (picture string) {

	containerPicture := h.DOM.Find("div.ut2-gl__image")

	picture, exists := containerPicture.ChildrenFiltered("a").ChildrenFiltered("img").Attr("data-src")
	if !exists {
		log.Printf("Wrong on get picture in: %s", h.Request.URL.String())
	}

	return

}
