package aversi

import (
	"fmt"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/model"
	"github.com/gocolly/colly"
)

const (
	urlresource = "https://www.aversi.ge"
	discountUrl = "/ka/aqciebi"
	place       = "Aversi"
)

type parserDiscount struct {
	Collector *colly.Collector
}

func NewParser() parsers.ParseDiscounter {

	return &parserDiscount{
		Collector: colly.NewCollector(),
	}
}

func (p *parserDiscount) ParseDiscounts() (*[]model.Discount, error) {

	var discounts []model.Discount

	// Need go on every page discount
	p.Collector.OnHTML("h5.entry-title.mt-0.pt-0", func(h *colly.HTMLElement) {
		h.Request.Visit(h.ChildAttr("a", "href"))
	})

	// Implement pagination
	p.Collector.OnHTML("ul.pagination.theme-colored", func(h *colly.HTMLElement) {

		listOfPagination := h.DOM.Find("li a")
		var isNextPage bool
		var nextPage string
		for _, node := range listOfPagination.Nodes {
			for _, atr := range node.Attr {
				if atr.Val == "next" {
					isNextPage = true
				}

				if atr.Key == "href" {
					nextPage = atr.Val
				}
			}
			if isNextPage {
				break
			}
		}

		if len(nextPage) > 0 {
			h.Request.Visit(nextPage)
		}

	})

	// Processing page discount
	p.Collector.OnHTML("article.post.clearfix.mb-50.pb-30", func(h *colly.HTMLElement) {

		discount := model.Discount{
			Url:     h.Request.URL.String(),
			Place:   place,
			Title:   h.DOM.Find(`h1.entry-title.font-22`).Text(),
			Picture: getPicture(h),
			//Description: getDescription(h),
		}

		discounts = append(discounts, discount)

	})

	p.Collector.Visit(fmt.Sprintf("%v%v", urlresource, discountUrl))

	return &discounts, nil
}

func getPicture(h *colly.HTMLElement) (picture string) {

	picture, exists := h.DOM.Find("img.img-fullwidth.img-responsive.aqciis-photo").Attr("src")
	if exists {
		picture = fmt.Sprintf("%v%v", urlresource, picture)
	}

	return

}

func getDescription(h *colly.HTMLElement) (description string) {

	description = h.DOM.Find("div.entry-content").Text()

	return

}
