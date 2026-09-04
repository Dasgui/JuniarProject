package handler

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/lib/json"
	"JuniarProject/internal/lib/parser"
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
	var parameters models.ProductsGetQueryParameters
	var err error

	var result []models.Product

	parameters, err = parseRequestParameters(r)
	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	result, err = h.service.GetProducts(r.Context(), parameters)

	if err != nil {
		internalErrors.PrintError(w, err)
		return
	}

	json.Write(w, http.StatusOK, result)
}

func parseRequestParameters(r *http.Request) (models.ProductsGetQueryParameters, error) {
	var result models.ProductsGetQueryParameters
	var err error

	result.Category = r.URL.Query().Get("category")

	result.PriceFrom, err = parser.ParseRequestToFloat(r, "price_from")
	if err != nil {
		return result, err
	}

	result.PriceTo, err = parser.ParseRequestToFloat(r, "price_to")
	if err != nil {
		return result, err
	}

	if result.PriceFrom != 0 && result.PriceTo != 0 && result.PriceFrom > result.PriceTo {
		return result, internalErrors.PriceError.Err
	}

	result.Limit, err = parser.ParseRequestToInt(r, "limit")
	if err != nil {
		return result, err
	}

	result.Offset, err = parser.ParseRequestToInt(r, "offset")
	if err != nil {
		return result, err
	}

	return result, nil
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
