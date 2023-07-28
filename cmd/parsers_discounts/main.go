package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/internal/parsers/aversi"
	"github.com/MaxBrainygame/Discounts-GE/internal/parsers/nikora"
	"github.com/MaxBrainygame/Discounts-GE/model"
)

func main() {

	parsersDiccount := GetParsersDiscount()
	// translatorDiscount := microsoftTranslate.NewTranslatorDiscount()

	for _, parser := range parsersDiccount {

		store, err := parser.ParseDiscounts()
		if err != nil {
			log.Fatal(err)
		}

		// discountsLanguage, err := translatorDiscount.TranslateDiscounts(discounts)
		// if err != nil {
		// 	log.Println(err)
		// }

		err = WriteDiscount(store, fmt.Sprintf("Discounts%s.json", store.Name))
		if err != nil {
			log.Fatal(err)
		}
		// err = WriteDiscount(discountsLanguage, "DiscountsLanguage.json")
		// if err != nil {
		// 	log.Fatal(err)
		// }

	}

}

func GetParsersDiscount() (parsersDiccount []parsers.ParseDiscounter) {

	parsersDiccount = append(parsersDiccount, aversi.NewParser())
	// parsersDiccount = append(parsersDiccount, cleanhouse.NewParser())
	// parsersDiccount = append(parsersDiccount, psp.NewParser())
	parsersDiccount = append(parsersDiccount, nikora.NewParser())

	return
}

func WriteDiscount[typeDiscount *model.Store | *[]model.DiscountLanguage](discounts typeDiscount,
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
