package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
)

type Store struct {
	ID        primitive.ObjectID `bson:"_id"`
	Logo      string             `bson:"logo"`
	Name      string             `bson:"name"`
	Host      string             `bson:"host"`
	Category  CategoryStores     `bson:"category"`
	Discounts []Discount         `bson:"discounts"`
}

type CategoryStores struct {
	Key  string `bson:"key"`
	Name string `bson:"name"`
}

type Discount struct {
	Url         string         `bson:"url"`
	Place       string         `bson:"place"`
	Picture     string         `bson:"picture"`
	Title       string         `bson:"title"`
	Description string         `bson:"description"`
	Goods       []DiscountItem `bson:"goods"`
}

type DiscountLanguage struct {
	Url         string
	Place       string
	Picture     string
	Title       string
	Description string
	Language    language.Tag
	Goods       []DiscountItem
}

type DiscountItem struct {
	Url          string  `bson:"url"`
	Picture      string  `bson:"picture"`
	Title        string  `bson:"title"`
	RegularPrice float64 `bson:"regular_price"`
	FinalPrice   float64 `bson:"final_price"`
}
