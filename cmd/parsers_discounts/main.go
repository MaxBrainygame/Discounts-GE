package main

import (
	"encoding/json"
	"fmt"
	"os"

	"log"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/internal/parsers/nikora"
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

	err = os.WriteFile(filename, discountsJson, 0644)
	if err != nil {
		return err
	}
	// 	ctx := context.Background()
	//
	// 	opts := options.Client().ApplyURI("mongodb://localhost:32771")
	//
	// 	client, err := mongo.Connect(ctx, opts)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	// Check the connection
	// 	err = client.Ping(ctx, nil)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	fmt.Println("Connected to MongoDB!")
	//
	// 	collection := client.Database("discounts").Collection("shops")
	//
	// 	_, err = collection.Indexes().CreateOne(
	// 		ctx,
	// 		mongo.IndexModel{
	// 			Keys: bson.D{
	// 				{Key: "host", Value: 1},
	// 				{Key: "name", Value: 1},
	// 			},
	// 			Options: options.Index().SetUnique(true),
	// 		},
	// 	)
	//
	// 	_, err = collection.InsertOne(context.Background(), discounts)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	return nil
}
