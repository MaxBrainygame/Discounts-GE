package main

import (
	"encoding/json"

	"io/ioutil"
	"log"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/internal/parsers/nikora"
	"github.com/MaxBrainygame/Discounts-GE/internal/translations/microsoftTranslate"
	"github.com/MaxBrainygame/Discounts-GE/model"
)

const (
	urlresource = "https://www.aversi.ge"
	discountUrl = "/ka/aqciebi"
)

func main() {

	parsersDiccount := GetParsersDiscount()
	translatorDiscount := microsoftTranslate.NewTranslatorDiscount()

	for _, parser := range parsersDiccount {

		discounts, err := parser.ParseDiscounts()
		if err != nil {
			log.Fatal(err)
		}

		discountsLanguage, err := translatorDiscount.TranslateDiscounts(discounts)
		if err != nil {
			log.Println(err)
		}

		err = WriteDiscount(discounts, "Discounts.json")
		if err != nil {
			log.Fatal(err)
		}
		err = WriteDiscount(discountsLanguage, "DiscountsLanguage.json")
		if err != nil {
			log.Fatal(err)
		}

	}

}

func GetParsersDiscount() (parsersDiccount []parsers.ParseDiscounter) {

	// parsersDiccount = append(parsersDiccount, aversi.NewParser())
	// parsersDiccount = append(parsersDiccount, cleanhouse.NewParser())
	// parsersDiccount = append(parsersDiccount, psp.NewParser())
	parsersDiccount = append(parsersDiccount, nikora.NewParser())

	return
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
