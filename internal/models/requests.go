package models

type ProductsGetQueryParameters struct {
	Category string

	PriceFrom float64
	PriceTo   float64

	Limit  int
	Offset int
}
