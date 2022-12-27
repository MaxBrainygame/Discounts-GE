package psp

type Categories struct {
	Data []Category
}

type Category struct {
	Id       int
	Name     string
	Level    int
	Url_path string
	Image    string
	Children []Category
}

type ListProducts struct {
	Success bool
	Errors  string
	Data    DataProducts
}

type DataProducts struct {
	Items     []ItemProduct
	Page_info PageInfo
}

type PageInfo struct {
	Current_page int
	Page_size    int
	Total_pages  int
}

type ItemProduct struct {
	Name        string
	Url_key     string
	Price_range PriceRange
	Thumbnail   Thumbnail
}

type PriceRange struct {
	Maximum_price MaximumPrice
}

type MaximumPrice struct {
	Final_price   FinalPrice
	Regular_price RegularPrice
}

type FinalPrice struct {
	Value float64
}

type RegularPrice struct {
	Value float64
}

type Thumbnail struct {
	Url string
}
