package aversi

import (
	"fmt"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/model"
	"github.com/gocolly/colly"
)

const (
	urlHost     = "https://www.aversi.ge"
	discountUrl = "/ka/aqciebi"
	place       = "Aversi"
	logoUrl     = "https://www.aversi.ge/images/logo.png"
	typeStore   = "Pharmacy"
)

type parserDiscount struct {
	Collector *colly.Collector
	Place     string
	LogoUrl   string
	Url       string
	UrlHost   string
	TypeStore string
}

func NewParser() parsers.ParseDiscounter {

	return &parserDiscount{
		Collector: colly.NewCollector(),
		Place:     place,
		LogoUrl:   logoUrl,
		Url:       fmt.Sprint(urlHost, discountUrl),
		UrlHost:   urlHost,
		TypeStore: typeStore,
	}
}

func (p *parserDiscount) ParseDiscounts(categoryStores map[string]*model.CategoryStores) (*model.Store, error) {

	var (
		discounts []model.Discount
		store     model.Store
	)

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
			Picture: p.getPicture(h),
		}

		discounts = append(discounts, discount)

	})

	p.Collector.Visit(p.Url)

	store = model.Store{
		Name:      p.Place,
		Logo:      p.LogoUrl,
		Host:      p.UrlHost,
		Category:  *categoryStores[p.TypeStore],
		Discounts: discounts,
	}

	return &store, nil
}

func (p *parserDiscount) getPicture(h *colly.HTMLElement) (picture string) {

	picture, exists := h.DOM.Find("img.img-fullwidth.img-responsive.aqciis-photo").Attr("src")
	if exists {
		picture = fmt.Sprintf("%v%v", p.UrlHost, p.UrlHost)
	}

	return

}
