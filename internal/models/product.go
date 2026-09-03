package models

import (
	"JuniarProject/internal/internalErrors"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Category    string           `json:"category"`
	CreatedAt   pgtype.Timestamp `json:"created_at" swaggerignore:"true"`
}

type ProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

func (r *ProductRequest) ConvertToProduct() Product {
	return Product{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Category:    r.Category,
	}
}

func (product *ProductRequest) CheckFields() error {
	if product.IsFieldsEmpty() {
		return internalErrors.EmptyFieldsError.Err
	}
	if product.IsNegativePrice() {
		return internalErrors.NegativePriceError.Err
	}
	return nil
}

func (product *ProductRequest) IsFieldsEmpty() bool {
	if strings.TrimSpace(product.Name) == "" || strings.TrimSpace(product.Description) == "" || strings.TrimSpace(product.Category) == "" {
		return true
	}
	return false
}

func (product *ProductRequest) IsNegativePrice() bool {
	if product.Price < 0 {
		return true
	}
	return false
}
