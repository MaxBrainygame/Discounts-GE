package myMemory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MaxBrainygame/Discounts-GE/model"
	"golang.org/x/text/language"
)

const (
	urlTranslate = "https://api.mymemory.translated.net/get"
)

type TranslatorDiscount struct {
	url             string
	defaultLanguage language.Tag
}

type Response struct {
	ResponseData ResponseData
}

type ResponseData struct {
	TranslatedText string
}

func NewTranslatorDiscount() (t *TranslatorDiscount) {
	return &TranslatorDiscount{
		url:             urlTranslate,
		defaultLanguage: language.Georgian,
	}
}

func (t *TranslatorDiscount) TranslateDiscounts(discounts *[]model.Discount) (*[]model.DiscountLanguage, error) {

	languages := [2]language.Tag{
		language.Russian,
		language.English,
	}

	var discountsLanguage []model.DiscountLanguage

	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		return nil, err
	}

	for _, lang := range languages {

		for _, discount := range *discounts {

			if len(discount.Title) == 0 {
				continue
			}

			translatedTitle, err := t.translatedText(req, discount.Title, &lang)
			if err != nil {
				return &discountsLanguage, err
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

	return &discountsLanguage, nil

}

func (t *TranslatorDiscount) translatedText(req *http.Request, forTrasnlate string, lang *language.Tag) (*Response, error) {

	var translated Response

	q := req.URL.Query()
	q.Set("q", forTrasnlate)
	q.Set("langpair", fmt.Sprintf("%v|%v", t.defaultLanguage.String(), lang.String()))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unsuccessful answer when tranlate. Status: %v", res.StatusCode)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &translated)
	if err != nil {
		return nil, err
	}

	return &translated, nil
}
