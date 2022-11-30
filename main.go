package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Discount struct {
	Place       string
	Title       string
	Description template.HTML
}

func main() {

	//ref := h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
	// log.Println(ref)
	// fmt.Println("Visiting", r.URL)
	discounts := ParseDiscounts()

	discountsJson, err := json.Marshal(discounts)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("Discounts.json", discountsJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func ParseDiscounts() (discounts []Discount) {

	collector := colly.NewCollector()

	// Need go on every page discount
	collector.OnHTML(`h5.entry-title.mt-0.pt-0`, func(h *colly.HTMLElement) {

		h.Request.Visit(h.ChildAttr("a", "href"))

	})

	// Processing page discount
	collector.OnHTML(`article.post.clearfix.mb-50.pb-30`, func(h *colly.HTMLElement) {

		discount := Discount{
			Place:       "Aversi",
			Title:       h.DOM.Find(`h1.entry-title.font-22`).Text(),
			Description: getDescription(h),
		}

		discounts = append(discounts, discount)

	})

	collector.Visit("https://www.aversi.ge/ka/aqciebi")

	return
}

func getDescription(h *colly.HTMLElement) (description template.HTML) {

	escapedAndJoined := template.HTMLEscaper(h.DOM.Html())
	description = template.HTML(strings.ReplaceAll(escapedAndJoined, "\n", "<br>"))
	// if err != nil {
	// 	log.Printf(`In parser 'Aversi' failed get Description. Err: %v`, err)
	// }

	return

}
