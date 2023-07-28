package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/MaxBrainygame/Discounts-GE/model"
)

func stores(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var (
		storeAversi model.Store
		storeNikora model.Store
		stores      []model.Store
	)
	file, err := ioutil.ReadFile("DiscountsAversi.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeAversi)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall")
	}

	file, err = ioutil.ReadFile("DiscountsNikora.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeNikora)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall. Err: %s", err)
	}

	stores = append(stores, storeAversi)
	stores = append(stores, storeNikora)

	resp, err := json.Marshal(stores)
	if err != nil {
		log.Fatalf("Error happened JSON Marshall. Err: %s", err)
	}

	w.Write(resp)
}

func promotions(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var (
		store    model.Store
		fileName string
	)

	storeHost := r.URL.Query().Get("store")
	if storeHost == "https://www.aversi.ge" {
		fileName = "DiscountsAversi.json"
	} else if storeHost == "http://nikorasupermarket.ge" {
		fileName = "DiscountsNikora.json"
	} else {
		w.WriteHeader(400)
		return
	}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &store)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall")
	}

	resp, err := json.Marshal(store.Discounts)
	if err != nil {
		log.Fatalf("Error happened JSON Marshall. Err: %s", err)
	}

	w.Write(resp)

}

func main() {
	http.HandleFunc("/stores", stores)
	http.HandleFunc("/promotions", promotions)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	// err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
