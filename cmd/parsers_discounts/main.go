package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/internal/parsers/psp"
	"github.com/MaxBrainygame/Discounts-GE/model"
)

func main() {

	categoryStores := make(map[string]*model.CategoryStores)
	categoryStores["Pharmacy"] = &model.CategoryStores{Key: "58270", Name: "Pharmacy"}
	categoryStores["Grocery"] = &model.CategoryStores{Key: "58780", Name: "Grocery"}
	categoryStores["CosmeticsHouseholdCleaningProducts"] = &model.CategoryStores{Key: "57702", Name: "Cosmetics & Hygiene"}

	parsersDiccount := GetParsersDiscount()
	// translatorDiscount := microsoftTranslate.NewTranslatorDiscount()

	for _, parser := range parsersDiccount {
		var s string
		store, err := parser.ParseDiscounts(categoryStores)
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

	// parsersDiccount = append(parsersDiccount, aversi.NewParser())
	// parsersDiccount = append(parsersDiccount, cleanhouse.NewParser())
	parsersDiccount = append(parsersDiccount, psp.NewParser())
	// parsersDiccount = append(parsersDiccount, nikora.NewParser())

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
