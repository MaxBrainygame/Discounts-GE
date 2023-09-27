package psp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/MaxBrainygame/Discounts-GE/internal/parsers"
	"github.com/MaxBrainygame/Discounts-GE/model"
	"github.com/gocolly/colly"
)

const (
	urlHost     = "https://psp.ge"
	productsUrl = "/category/categoryName/products"
	categoryUrl = "/category/tree"
	place       = "PSP"
	logoUrl     = "https://psp.ge/logo.png"
	typeStore   = "Pharmacy"
	// Quantity elements in one page
	itemsPerPage = "600"
)

type parserDiscount struct {
	Collector    *colly.Collector
	Place        string
	LogoUrl      string
	UrlHost      string
	ProductsUrl  string
	CategoryUrl  string
	ItemsPerPage string
	TypeStore    string
}

func NewParser() parsers.ParseDiscounter {

	return &parserDiscount{
		Collector:    colly.NewCollector(),
		Place:        place,
		LogoUrl:      logoUrl,
		UrlHost:      urlHost,
		ProductsUrl:  fmt.Sprint(urlHost, productsUrl),
		CategoryUrl:  fmt.Sprint(urlHost, categoryUrl),
		ItemsPerPage: itemsPerPage,
		TypeStore:    typeStore,
	}
}

func (p *parserDiscount) ParseDiscounts(categoryStores map[string]*model.CategoryStores) (*model.Store, error) {

	var (
		discounts []model.Discount
		discount  model.Discount
		store     model.Store
	)

	treeCategories, err := p.getTreeCategories()
	if err != nil {
		return &store, nil
	}

	store = model.Store{
		Name:     p.Place,
		Logo:     p.LogoUrl,
		Host:     p.UrlHost,
		Category: *categoryStores[p.TypeStore],
	}

	// Starting caterories with Level = 2
	for _, category := range treeCategories.Data[0].Children {

		discount = model.Discount{
			Url:     fmt.Sprintf("%s/%s.html", p.UrlHost, category.Url_path),
			Place:   p.Place,
			Title:   category.Name,
			Picture: category.Image,
		}

		currentPage := 1
		totalPages := 0
		for true {

			listProduct, err := p.getProducts(currentPage, category.Id)
			if err != nil {
				store.Discounts = discounts
				return &store, err
			}

			for _, product := range listProduct.Data.Items {

				regularPrice := product.Price_range.Maximum_price.Regular_price.Value
				finalPrice := product.Price_range.Maximum_price.Final_price.Value
				if finalPrice >= regularPrice {
					continue
				}

				discountItem := model.DiscountItem{
					Url:          fmt.Sprintf("%s/%s.html", p.UrlHost, product.Url_key),
					Picture:      product.Thumbnail.Url,
					Title:        product.Name,
					RegularPrice: regularPrice,
					FinalPrice:   finalPrice,
				}

				discount.Goods = append(discount.Goods, discountItem)

			}

			if totalPages == 0 {
				totalPages = listProduct.Data.Page_info.Total_pages
			}

			if currentPage == totalPages {
				break
			}

			currentPage++

		}

		discounts = append(discounts, discount)

	}

	store.Discounts = discounts

	return &store, nil

}

func (p *parserDiscount) getProducts(currentPage int, categoryId int) (*ListProducts, error) {

	var listProduct ListProducts

	urlPage := strings.ReplaceAll(p.ProductsUrl, "categoryName", strconv.Itoa(categoryId))
	req, err := http.NewRequest(http.MethodGet, urlPage, nil)
	if err != nil {
		return nil, err
	}

	reqParam := req.URL.Query()
	reqParam.Set("currentPage", strconv.Itoa(currentPage))
	reqParam.Set("itemsPerPage", p.ItemsPerPage)

	req.URL.RawQuery = reqParam.Encode()

	contentBody, err := p.executeRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*contentBody, &listProduct)
	if err != nil {
		return nil, err
	}

	if !listProduct.Success {
		err = errors.New(listProduct.Errors)
		return nil, err
	}

	return &listProduct, nil

}

func (p *parserDiscount) getTreeCategories() (*Categories, error) {

	var categories Categories

	req, err := http.NewRequest(http.MethodGet, p.CategoryUrl, nil)
	if err != nil {
		return nil, err
	}

	contentBody, err := p.executeRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*contentBody, &categories)
	if err != nil {
		return nil, err
	}

	return &categories, nil

}

func (p *parserDiscount) executeRequest(req *http.Request) (*[]byte, error) {

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unsuccessful answer when tranlate. Status: %v", resp.StatusCode)
		return nil, err
	}

	defer resp.Body.Close()
	contentBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &contentBody, nil

}
