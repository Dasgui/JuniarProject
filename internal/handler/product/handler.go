package handler

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/lib/json"
	"JuniarProject/internal/models"
	service "JuniarProject/internal/service/product"
	"errors"
	"log"
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
	var priceFrom float64
	var priceTo float64

	var err error

	limit := 0
	offset := 0

	category := r.URL.Query().Get("category")

	if priceFromStr := r.URL.Query().Get("price_from"); priceFromStr != "" {
		if val, err := strconv.ParseFloat(priceFromStr, 64); err == nil {
			priceFrom = val
		} else {
			log.Println(err)
			http.Error(w, internalErrors.ParametresRequestError.Error(), http.StatusBadRequest)
			return
		}
	}

	if priceToStr := r.URL.Query().Get("price_to"); priceToStr != "" {
		if val, err := strconv.ParseFloat(priceToStr, 64); err == nil {
			priceTo = val
		} else {
			log.Println(err)
			http.Error(w, internalErrors.ParametresRequestError.Error(), http.StatusBadRequest)
			return
		}
	}

	// Валидация: price_from не может быть больше price_to
	if priceFrom != 0 && priceTo != 0 && priceFrom > priceTo {
		http.Error(w, internalErrors.PriceError.Error(), http.StatusBadRequest)
		return
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err != nil || val < 0 {
			log.Println(err)
			http.Error(w, internalErrors.ParametresRequestError.Error(), http.StatusBadRequest)
			return
		} else {
			limit = val
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err != nil || val < 0 {
			log.Println(err)
			http.Error(w, internalErrors.ParametresRequestError.Error(), http.StatusBadRequest)
			return
		} else {
			offset = val
		}
	}

	result, err := h.service.GetProducts(r.Context(), category, priceFrom, priceTo, limit, offset)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
		http.Error(w, internalErrors.IdRequestError.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		log.Println(err)
		http.Error(w, internalErrors.IdRequestError.Error(), http.StatusBadRequest)
		return
	}

	var result models.Product
	if err := json.Read(r, &result); err != nil {
		log.Println(err)
		http.Error(w, internalErrors.JsonTypeError.Error(), http.StatusBadRequest)
		return
	}

	result, err = h.service.UpdateProduct(r.Context(), result, id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		log.Println(err)
		http.Error(w, internalErrors.IdRequestError.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusOK, result)
}
