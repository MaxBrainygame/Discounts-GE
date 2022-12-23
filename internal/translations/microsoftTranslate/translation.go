package microsoftTranslate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MaxBrainygame/Discounts-GE/model"
	"golang.org/x/text/language"
)

const (
	urlTranslate = "https://microsoft-translator-text.p.rapidapi.com/translate"
)

type TranslatorDiscount struct {
	url             string
	defaultLanguage language.Tag
	languages       [2]language.Tag
}

type ForTranslate struct {
	Text string
}

type ResponseTranslate struct {
	TranslatedMicrosoft []TranslatedMicrosoft
}

type TranslatedMicrosoft struct {
	Translations []Translations
}

type Translations struct {
	Text string
	To   string
}

func NewTranslatorDiscount() (t *TranslatorDiscount) {
	return &TranslatorDiscount{
		url:             urlTranslate,
		defaultLanguage: language.Georgian,
		languages:       [2]language.Tag{language.English, language.Russian},
	}
}

func (t *TranslatorDiscount) TranslateDiscounts(discounts *[]model.Discount) (*[]model.DiscountLanguage, error) {

	var discountsLanguage []model.DiscountLanguage
	var discountLanguage model.DiscountLanguage

	for _, discount := range *discounts {

		if len(discount.Title) == 0 {
			continue
		}

		translatedTexts, err := t.translatedText(discount.Title, "")
		if err != nil {
			return &discountsLanguage, err
		}

		for _, translation := range translatedTexts[0].Translations {

			discountLanguage = model.DiscountLanguage{
				Url:      discount.Url,
				Place:    discount.Place,
				Picture:  discount.Picture,
				Title:    translation.Text,
				Language: language.Make(translation.To),
			}

			for _, discountItem := range discount.Goods {

				if len(discountItem.Title) == 0 {
					continue
				}

				translatedTexts, err := t.translatedText(discountItem.Title, translation.To)
				if err != nil {
					return &discountsLanguage, err
				}

				for _, translationItem := range translatedTexts[0].Translations {

					discountItemLanguage := model.DiscountItem{
						Url:          discountItem.Url,
						Picture:      discountItem.Picture,
						Title:        translationItem.Text,
						RegularPrice: discountItem.RegularPrice,
						FinalPrice:   discountItem.FinalPrice,
					}

					discountLanguage.Goods = append(discountLanguage.Goods, discountItemLanguage)

				}

			}

			discountsLanguage = append(discountsLanguage, discountLanguage)

			time.Sleep(2 * time.Second)

		}

	}

	return &discountsLanguage, nil
}

func (t *TranslatorDiscount) translatedText(forTrasnlateText string, lang string) ([]TranslatedMicrosoft, error) {

	var translated []TranslatedMicrosoft
	var forTrasnlate [1]ForTranslate

	forTrasnlate[0] = ForTranslate{Text: forTrasnlateText}

	reqBody, err := json.Marshal(&forTrasnlate)
	if err != nil {
		return nil, err
	}

	// reqBody := []byte(fmt.Sprintf(`[
	// 	{
	// 		"Text": "%s"
	// 	}
	// ]`, forTrasnlateText))

	// reqBody := []byte(`[
	// 	{
	// 		"Text": "BELITA  შამპუნი თხელი და სუსტი თმისთვის 480 მლ (ბელიტა)"
	// 	}
	// ]`)

	req, err := http.NewRequest(http.MethodPost, t.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-RapidAPI-Key", "08c23b896emsh168bd72a2444c26p167e4ajsnf68a3d859a27")
	req.Header.Set("X-RapidAPI-Host", "microsoft-translator-text.p.rapidapi.com")
	req.Header.Set("Content-Type", "application/json")

	if len(lang) == 0 {
		lang = fmt.Sprintf("%s,%s", t.languages[0].String(), t.languages[1].String())
	}

	q := req.URL.Query()
	q.Set("api-version", "3.0")
	q.Set("profanityAction", "NoAction")
	q.Set("textType", "plain")
	q.Set("from", t.defaultLanguage.String())
	q.Set("to", lang)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unsuccessful answer tranlate. Status: %v", res.Status)
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

	return translated, nil
}
