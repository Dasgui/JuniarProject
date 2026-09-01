package handler

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/lib/json"
	"JuniarProject/internal/models"
	service "JuniarProject/internal/service/product"
	"errors"
	"log"
	"net/http"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var result models.Product

	if err := json.Read(r, &result); err != nil {
		log.Println(err)
		http.Error(w, internalErrors.JsonTypeError.Error(), http.StatusBadRequest)
		return
	}

	switch err := result.CheckFields(); {
	case errors.Is(err, internalErrors.NegativePriceError):
		http.Error(w, internalErrors.NegativePriceError.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, internalErrors.EmptyFieldsErr):
		http.Error(w, internalErrors.EmptyFieldsErr.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.CreatProduct(r.Context(), result)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetProducts(r.Context())

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusOK, result)
}
