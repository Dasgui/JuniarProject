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
	CreatedAt   pgtype.Timestamp `json:"created_at"`
}

func (product *Product) CheckFields() error {
	if product.IsFieldsEmpty() {
		return internalErrors.EmptyFieldsError
	}
	if product.IsNegativePrice() {
		return internalErrors.NegativePriceError
	}
	return nil
}

func (product *Product) IsFieldsEmpty() bool {
	if strings.TrimSpace(product.Name) == "" || strings.TrimSpace(product.Description) == "" || strings.TrimSpace(product.Category) == "" {
		return true
	}
	return false
}

func (product *Product) IsNegativePrice() bool {
	if product.Price < 0 {
		return true
	}
	return false
}
