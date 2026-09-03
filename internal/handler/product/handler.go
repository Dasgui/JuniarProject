package handler

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/lib/json"
	"JuniarProject/internal/models"
	service "JuniarProject/internal/service/product"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		internalErrors.PrintError(w, err)
		return
	}

	err := result.CheckFields()
	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	result, err = h.service.CreateProduct(r.Context(), result)

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	var priceFrom float64
	var priceTo float64

	var err error

	limit := 0
	offset := 0

	category := r.URL.Query().Get("category")

	// Необходимо будет сделать общую функцию
	if priceFromStr := r.URL.Query().Get("price_from"); priceFromStr != "" {
		if val, err := strconv.ParseFloat(priceFromStr, 64); err == nil {
			if val < 0 {
				internalErrors.PrintError(w, internalErrors.NegativePriceError.Err)
				return
			}
			priceFrom = val
		} else {
			internalErrors.PrintError(w, err)
			return
		}
	}

	if priceToStr := r.URL.Query().Get("price_to"); priceToStr != "" {
		if val, err := strconv.ParseFloat(priceToStr, 64); err == nil {
			if val < 0 {
				internalErrors.PrintError(w, internalErrors.NegativePriceError.Err)
				return
			}
			priceTo = val
		} else {
			internalErrors.PrintError(w, err)
			return
		}
	}

	// Валидация: price_from не может быть больше price_to
	if priceFrom != 0 && priceTo != 0 && priceFrom > priceTo {
		internalErrors.PrintError(w, internalErrors.PriceError.Err)
		return
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err != nil || val < 0 {
			internalErrors.PrintError(w, err)
			return
		} else {
			limit = val
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err != nil || val < 0 {
			internalErrors.PrintError(w, err)
			return
		} else {
			offset = val
		}
	}

	result, err := h.service.GetProducts(r.Context(), category, priceFrom, priceTo, limit, offset)

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	result, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {

		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	var result models.Product
	if err := json.Read(r, &result); err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	result, err = h.service.UpdateProduct(r.Context(), result, id)
	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	result, err := h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}
