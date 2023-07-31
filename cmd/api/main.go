package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

func products(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var (
		store        model.Store
		findDiscount model.Discount
		fileName     string
	)

	discountHost := r.URL.Query().Get("discount")

	if strings.Contains(discountHost, "https://www.aversi.ge") {
		fileName = "DiscountsAversi.json"
	} else if strings.Contains(discountHost, "http://nikorasupermarket.ge") {
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

	for _, discount := range store.Discounts {

		if discount.Url == discountHost {
			findDiscount = discount
			break
		}

	}

	if findDiscount.Url == "" {
		w.WriteHeader(400)
		return
	}

	resp, err := json.Marshal(findDiscount.Goods)
	if err != nil {
		log.Fatalf("Error happened JSON Marshall. Err: %s", err)
	}

	w.Write(resp)

}

func main() {

	server := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}
	}()

	http.HandleFunc("/stores", stores)
	http.HandleFunc("/promotions", promotions)
	http.HandleFunc("/products", products)
	err := server.ListenAndServe()
	// err := http.ListenAndServe(":8080", nil)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}

func closeApp() {

}
