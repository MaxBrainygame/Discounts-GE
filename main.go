package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/MaxBrainygame/Discounts-GE/model"
	"github.com/gocolly/colly"
	"golang.org/x/text/language"
)

const (
	urlresource = "https://www.aversi.ge"
	discountUrl = "/ka/aqciebi"
)

type Response struct {
	ResponseData ResponseData
}

type ResponseData struct {
	TranslatedText string
}

func main() {

	discounts := ParseDiscounts()

	discountsLanguage := TranslateDiscounts(&discounts)

	err := WriteDiscount(&discounts, "Discounts.json")
	if err != nil {
		log.Fatal(err)
	}
	err = WriteDiscount(&discountsLanguage, "DiscountsLanguage.json")
	if err != nil {
		log.Fatal(err)
	}

}

func WriteDiscount[typeDiscount *[]model.Discount | *[]model.DiscountLanguage](discounts typeDiscount,
	filename string) error {

	discountsJson, err := json.Marshal(discounts)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, discountsJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ParseDiscounts() (discounts []model.Discount) {

	collector := colly.NewCollector()

	// Need go on every page discount
	collector.OnHTML(`h5.entry-title.mt-0.pt-0`, func(h *colly.HTMLElement) {

		h.Request.Visit(h.ChildAttr("a", "href"))

	})

	// Processing page discount
	collector.OnHTML(`article.post.clearfix.mb-50.pb-30`, func(h *colly.HTMLElement) {

		discount := model.Discount{
			Url:     h.Request.URL.String(),
			Place:   "Aversi",
			Title:   h.DOM.Find(`h1.entry-title.font-22`).Text(),
			Picture: getPicture(h),
			//Description: getDescription(h),
		}

		discounts = append(discounts, discount)

	})

	collector.Visit(fmt.Sprintf("%v%v", urlresource, discountUrl))

	return
}

func TranslateDiscounts(discounts *[]model.Discount) (discountsLanguage []model.DiscountLanguage) {

	languages := [2]language.Tag{
		language.Russian,
		language.English,
	}

	urlTranslate := "https://api.mymemory.translated.net/get"

	req, err := http.NewRequest("GET", urlTranslate, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, lang := range languages {

		for _, discount := range *discounts {

			translatedTitle, err := translatedText(req, &discount.Title, &lang)
			if err != nil {
				fmt.Println(err)
				return
			}

			discountLanguage := model.DiscountLanguage{
				Url:      discount.Url,
				Place:    discount.Place,
				Picture:  discount.Picture,
				Title:    translatedTitle.ResponseData.TranslatedText,
				Language: lang,
			}

			discountsLanguage = append(discountsLanguage, discountLanguage)
		}

	}

	return

}

func translatedText(req *http.Request, forTrasnlate *string, lang *language.Tag) (translated Response, err error) {

	q := req.URL.Query()
	q.Set("q", *forTrasnlate)
	q.Set("langpair", fmt.Sprintf("%v|%v", language.Georgian.String(), lang.String()))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &translated)
	if err != nil {
		return
	}

	return
}

func getDescription(h *colly.HTMLElement) (description string) {

	description = h.DOM.Find("div.entry-content").Text()

	return

}

func getPicture(h *colly.HTMLElement) (picture string) {

	picture = h.DOM.Find("div.entry-content").Text()

	picture, exists := h.DOM.Find("img.img-fullwidth.img-responsive.aqciis-photo").Attr("src")
	if exists {
		picture = fmt.Sprintf("%v%v", urlresource, picture)
	}

	return

}
