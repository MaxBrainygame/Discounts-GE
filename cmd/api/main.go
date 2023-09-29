package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MaxBrainygame/Discounts-GE/model"
)

func categories(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	langCode := r.URL.Query().Get("lang")
	fmt.Println(langCode)
	var (
		storeAversi     model.Store
		storeNikora     model.Store
		storePSP        model.Store
		storeCleanHouse model.Store
		categoriesStore []model.CategoryStores
	)
	file, err := os.ReadFile("DiscountsAversi.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeAversi)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall")
	}

	file, err = os.ReadFile("DiscountsNikora.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeNikora)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall. Err: %s", err)
	}

	file, err = os.ReadFile("DiscountsPSP.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storePSP)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall. Err: %s", err)
	}

	file, err = os.ReadFile("DiscountsCleanHouse.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeCleanHouse)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall. Err: %s", err)
	}

	categoriesStore = append(categoriesStore, storeAversi.Category)
	categoriesStore = append(categoriesStore, storeNikora.Category)
	categoriesStore = append(categoriesStore, storePSP.Category)
	categoriesStore = append(categoriesStore, storeCleanHouse.Category)

	resp, err := json.Marshal(categoriesStore)
	if err != nil {
		log.Fatalf("Error happened JSON Marshall. Err: %s", err)
	}

	w.Write(resp)
}

func stores(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var (
		storeAversi     model.Store
		storeNikora     model.Store
		storePSP        model.Store
		storeCleanHouse model.Store
		stores          []model.Store
	)

	keyCategory := r.URL.Query().Get("key")

	file, err := os.ReadFile("DiscountsAversi.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeAversi)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall")
	}

	file, err = os.ReadFile("DiscountsNikora.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeNikora)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall. Err: %s", err)
	}

	file, err = os.ReadFile("DiscountsPSP.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storePSP)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall. Err: %s", err)
	}

	file, err = os.ReadFile("DiscountsCleanHouse.json")
	if err != nil {
		log.Fatalf("Error happened read file. Err: %s", err)
	}
	err = json.Unmarshal(file, &storeCleanHouse)
	if err != nil {
		log.Fatalf("Error happened JSON unmarhall. Err: %s", err)
	}

	stores = append(stores, storeAversi)
	stores = append(stores, storeNikora)
	stores = append(stores, storePSP)
	stores = append(stores, storeCleanHouse)

	result := []model.Store{}

	for i := 0; i < len(stores); i++ {

		if stores[i].Category.Key != keyCategory {
			continue
		}

		// stores = append(stores[:i], stores[i+1:]...)
		result = append(result, stores[i])

	}
	resp, err := json.Marshal(result)
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
	} else if storeHost == "https://psp.ge" {
		fileName = "DiscountsPSP.json"
	} else if storeHost == "https://ch.ge" {
		fileName = "DiscountsCleanHouse.json"
	} else {
		w.WriteHeader(400)
		return
	}

	file, err := os.ReadFile(fileName)
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
	} else if strings.Contains(discountHost, "https://psp.ge") {
		fileName = "DiscountsPSP.json"
	} else if strings.Contains(discountHost, "https://ch.ge") {
		fileName = "DiscountsCleanHouse.json"
	} else {
		w.WriteHeader(400)
		return
	}

	file, err := os.ReadFile(fileName)
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

	// 	cfg := config.MustLoad()
	//
	// 	log := logger.SetupLogger(cfg.Env)
	//
	// 	log.Info("starting API - discounts_ge", slog.String("env", cfg.Env))
	// 	log.Debug("debug enable")

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
	http.HandleFunc("/categories", categories)
	err := server.ListenAndServe()
	// err := http.ListenAndServe(":8080", nil)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
