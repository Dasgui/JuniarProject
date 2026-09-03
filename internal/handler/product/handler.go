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

// @Summary Добавить новый продукт
// @Tags products
// @Description Позволяет добавить новый продукт в базу данных
// @Accept json
// @Produce json
// @Param input body models.ProductRequest true "Информация о товаре"
// @Success 200 {object} models.Product
// @Failure 400 {string} internalErrors.InvalidParameterError
// @Failure 500 {string} internalErrors.InternalServerErr
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var requestBody models.ProductRequest
	var result models.Product

	if err := json.Read(r, &requestBody); err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	err := requestBody.CheckFields()
	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	result = requestBody.ConvertToProduct()

	result, err = h.service.CreateProduct(r.Context(), result)

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

// @Summary Получить все продукты
// @Tags products
// @Description Позволяет получить список продуктов, с возможностью фильтрации.
// @Accept json
// @Produce json
// @Param category query string false "Фильтр по категории"
// @Param price_from query number false "Минимальная цена"
// @Param price_to query number false "Максимальная цена"
// @Param limit query int false "Количество записей"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.Product
// @Failure 400 {string} internalErrors.InvalidParameterError
// @Failure 500 {string} internalErrors.InternalServerErr
// @Router /products [get]
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
	//var response []models.ProductResponse
	//for _, p := range result {
	//	response = append(response, models.ProductResponse{
	//		Id:        p.Id,
	//		Name:      p.Name,
	//		Category:  p.Category,
	//		Price:     p.Price,
	//		CreatedAt: p.CreatedAt.Time,
	//	})
	//}

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

// @Summary Получить конкретный продукт
// @Tags products
// @Description Позволяет получить конкретный продукт по его id
// @Accept json
// @Produce json
// @Param id path string false "Id товара"
// @Success 200 {object} models.Product
// @Failure 400 {string} internalErrors.InvalidParameterError
// @Failure 404 {string} internalErrors.DataNotFound
// @Failure 500 {string} internalErrors.InternalServerErr
// @Router /products/{id} [get]
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

// @Summary Обновить конкретный продукт
// @Tags products
// @Description Позволяет обновить конкретный продукт по его id
// @Accept json
// @Produce json
// @Param id path string false "Id товара"
// @Param input body models.ProductRequest true "Информация о товаре, которую необходимо обновить"
// @Success 200 {object} models.Product
// @Failure 400 {string} internalErrors.InvalidParameterError
// @Failure 404 {string} internalErrors.DataNotFound
// @Failure 500 {string} internalErrors.InternalServerErr
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	var requestBody models.ProductRequest
	var result models.Product

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	if err := json.Read(r, &requestBody); err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	result, err = h.service.UpdateProduct(r.Context(), requestBody, id)
	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

// @Summary Удалить конкретный продукт
// @Tags products
// @Description Позволяет удалить конкретный продукт по его id
// @Accept json
// @Produce json
// @Param id path string false "Id товара"
// @Success 200 {object} models.Product
// @Failure 400 {string} internalErrors.InvalidParameterError
// @Failure 404 {string} internalErrors.DataNotFound
// @Failure 500 {string} internalErrors.InternalServerErr
// @Router /products/{id} [delete]
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
