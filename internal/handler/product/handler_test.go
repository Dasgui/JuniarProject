package handler

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/lib/test_data"
	"JuniarProject/internal/models"
	mock_service "JuniarProject/internal/service/product/mocks"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestGet(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2026-01-01 12:00:00")
	expectedProducts := []models.Product{
		{
			Id:          1,
			Name:        "Test-1",
			Description: "Description Test-1",
			Price:       2300,
			Category:    "Middle",
			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
		},
		{
			Id:          2,
			Name:        "Test-2",
			Description: "Description Test-2",
			Price:       6500,
			Category:    "Electro",
			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
		},
		{
			Id:          3,
			Name:        "Test-1",
			Description: "Description Test-1",
			Price:       8300,
			Category:    "Electro",
			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockProductService(ctrl)
	h := NewProductHandler(service)

	r := chi.NewRouter()
	r.Get("/products/", h.GetProducts)

	successTestTable := []test_data.ProductTestTable{
		{
			Name:        "OK_NO_FILTER",
			QueryParams: "",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{}).
					Return(expectedProducts, nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  len(expectedProducts),
		},

		{
			Name:        "OK_FILTER_CATEGORY",
			QueryParams: "?category=Electro",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						Category: "Electro",
					}).
					Return(expectedProducts[:2], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  2,
		},

		{
			Name:        "OK_FILTER_PRICE_FROM",
			QueryParams: "?price_from=6000",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						PriceFrom: 6000.0,
					}).
					Return(expectedProducts[1:], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  2,
		},

		{
			Name:        "OK_FILTER_PRICE_TO",
			QueryParams: "?price_to=7000",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						PriceTo: 7000.0,
					}).
					Return(expectedProducts[:2], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  2,
		},

		{
			Name:        "OK_FILTER_PRICE_FROM_TO",
			QueryParams: "?price_from=6000&price_to=7000",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						PriceFrom: 6000.0, PriceTo: 7000.0,
					}).
					Return(expectedProducts[2:3], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  1,
		},

		{
			Name:        "OK_FILTER_CATEGORY_PRICE_FROM_TO",
			QueryParams: "?category=Middle&price_from=2000&price_to=10000",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						Category: "Middle", PriceFrom: 2000.0, PriceTo: 10000.0,
					}).
					Return(expectedProducts[1:2], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  1,
		},

		{
			Name:        "OK_FILTER_LIMIT",
			QueryParams: "?limit=1",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						Limit: 1,
					}).
					Return(expectedProducts[1:2], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  1,
		},

		{
			Name:        "OK_FILTER_OFFSET",
			QueryParams: "?offset=2",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						Offset: 2,
					}).
					Return(expectedProducts[:2], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  2,
		},

		{
			Name:        "OK_FILTER_LIMIT_OFFSET",
			QueryParams: "?limit=2&&offset=1",
			Setup: func() {
				service.EXPECT().
					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
						Limit: 2, Offset: 1,
					}).
					Return(expectedProducts[1:3], nil)
			},

			Except:         expectedProducts,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
			ExpectedCount:  2,
		},
	}
	errorTestTable := []test_data.ProductTestTable{
		{
			Name:        "ERROR_FILTER_SYNTEX_PRICE_FROM",
			QueryParams: "?price_from=asdasdas",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_NEGATIVE_PRICE_FROM",
			QueryParams: "?price_from=-200",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.NegativePriceError.Code,
			ExpectedError:  internalErrors.NegativePriceError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_SYNTEX_PRICE_TO",
			QueryParams: "?price_to=asdasd",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_NEGATIVE_PRICE_TO",
			QueryParams: "?price_to=-200",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.NegativePriceError.Code,
			ExpectedError:  internalErrors.NegativePriceError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_SYNTEX_PRICE_FROM_TO",
			QueryParams: "?price_from=asd&price_to=asdasd",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_NEGATIVE_PRICE_FROM_TO",
			QueryParams: "?price_from=-200&price_to=-200",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.NegativePriceError.Code,
			ExpectedError:  internalErrors.NegativePriceError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_PRICE_SYNTEX_FROM_NEGATIVE_TO",
			QueryParams: "?price_from=asd&price_to=-200",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_PRICE_NEGATIVE_FROM_SYNTEX_TO",
			QueryParams: "?price_from=-200&price_to=asd",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.NegativePriceError.Code,
			ExpectedError:  internalErrors.NegativePriceError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_SYNTEX_LIMIT",
			QueryParams: "?limit=asd",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_SYNTEX_OFFSET",
			QueryParams: "?offset=asd",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
			ExpectedCount:  0,
		},

		{
			Name:        "ERROR_FILTER_SYNTEX_LIMIT_OFFSET",
			QueryParams: "?limit=asd&?offset=asd",
			Setup:       func() {},

			Except:         expectedProducts,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
			ExpectedCount:  0,
		},
	}

	test_data.StartTest(t, successTestTable, r, http.MethodGet)
	test_data.StartTest(t, errorTestTable, r, http.MethodGet)
}

func TestGetById(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2026-01-01 12:00:00")
	expectedProduct := models.Product{
		Id:          1,
		Name:        "Test-1",
		Description: "Description Test-1",
		Price:       2300,
		Category:    "Middle",
		CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockProductService(ctrl)
	h := NewProductHandler(service)

	r := chi.NewRouter()
	r.Get("/products/{id}", h.GetProductByID)

	successTestTable := []test_data.ProductTestTable{
		{
			Name:        "OK",
			QueryParams: "1",
			Setup: func() {
				service.EXPECT().
					GetProductByID(gomock.Any(), 1).
					Return(expectedProduct, nil)
			},

			Except:         expectedProduct,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
		},
	}
	errorsIdTestTable := []test_data.ProductTestTable{
		{
			Name:        "ERROR_INVALID_ID",
			QueryParams: "asdasd",
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
		},
		{
			Name:        "ERROR_RANGE_ID",
			QueryParams: "200000000000000000000000000000000000000000000000000000000000000000000",
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.IdRangeError.Code,
			ExpectedError:  internalErrors.IdRangeError.Err,
		},

		{
			Name:        "ERROR_NOT_FOUND_ID",
			QueryParams: "200",
			Setup: func() {
				service.EXPECT().
					GetProductByID(gomock.Any(), 200).
					Return(models.Product{}, pgx.ErrNoRows)
			},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.DataNotFound.Code,
			ExpectedError:  internalErrors.DataNotFound.Err,
		},
	}

	test_data.StartTest(t, successTestTable, r, http.MethodGet)
	test_data.StartTest(t, errorsIdTestTable, r, http.MethodGet)
}

func TestDelete(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2026-01-01 12:00:00")
	expectedProduct := models.Product{
		Id:          1,
		Name:        "Test-1",
		Description: "Description Test-1",
		Price:       2300,
		Category:    "Middle",
		CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockProductService(ctrl)
	h := NewProductHandler(service)

	r := chi.NewRouter()
	r.Delete("/products/{id}", h.DeleteProduct)

	successTestTable := []test_data.ProductTestTable{
		{
			Name:        "OK",
			QueryParams: "1",
			Setup: func() {
				service.EXPECT().
					DeleteProduct(gomock.Any(), 1).
					Return(expectedProduct, nil)
			},

			Except:         expectedProduct,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
		},
	}

	errorsIdTestTable := []test_data.ProductTestTable{
		{
			Name:        "ERROR_INVALID_ID",
			QueryParams: "asdasd",
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.InvalidParameterError.Code,
			ExpectedError:  internalErrors.InvalidParameterError.Err,
		},

		{
			Name:        "ERROR_NOT_FOUND_ID",
			QueryParams: "200",
			Setup: func() {
				service.EXPECT().
					DeleteProduct(gomock.Any(), 200).
					Return(models.Product{}, pgx.ErrNoRows)
			},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.DataNotFound.Code,
			ExpectedError:  internalErrors.DataNotFound.Err,
		},

		{
			Name:        "ERROR_RANGE_ID",
			QueryParams: "200000000000000000000000000000000000000000000000000000000000000000000",
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.IdRangeError.Code,
			ExpectedError:  internalErrors.IdRangeError.Err,
		},
	}

	test_data.StartTest(t, successTestTable, r, http.MethodDelete)
	test_data.StartTest(t, errorsIdTestTable, r, http.MethodDelete)
}

func TestCreate(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2026-01-01 12:00:00")
	expectedProduct := models.Product{
		Id:          1,
		Name:        "Test-1",
		Description: "Description Test-1",
		Price:       2300,
		Category:    "Middle",
		CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockProductService(ctrl)
	h := NewProductHandler(service)

	r := chi.NewRouter()
	r.Post("/products/", h.CreateProduct)

	successTestTable := []test_data.ProductTestTable{
		{
			Name:        "OK",
			QueryParams: "",
			Body:        `{"name":"Test-1", "description": "Description Test-1", "price": 2300, "category": "Middle"}`,
			Setup: func() {
				service.EXPECT().
					CreateProduct(gomock.Any(), gomock.Any()).
					Return(expectedProduct, nil)
			},

			Except:         expectedProduct,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
		},
	}

	errorsEmptyFieldTestTable := []test_data.ProductTestTable{
		{
			Name:           "ERROR_EMPTY_NAME",
			Body:           `{"name":"", "description": "Description Test-1", "price": 2300, "category": "Middle"}`,
			Setup:          func() {},
			Except:         models.Product{},
			ExpectedStatus: internalErrors.EmptyFieldsError.Code,
			ExpectedError:  internalErrors.EmptyFieldsError.Err,
		},

		{
			Name:           "ERROR_EMPTY_TEXT_FIELDS",
			Body:           `{"name":"", "description": "", "price": 2300, "category": ""}`,
			Setup:          func() {},
			Except:         models.Product{},
			ExpectedStatus: internalErrors.EmptyFieldsError.Code,
			ExpectedError:  internalErrors.EmptyFieldsError.Err,
		},

		{
			Name:           "ERROR_EMPTY_BODY",
			Body:           `{}`,
			Setup:          func() {},
			Except:         models.Product{},
			ExpectedStatus: internalErrors.EmptyFieldsError.Code,
			ExpectedError:  internalErrors.EmptyFieldsError.Err,
		},
	}
	errorsFieldTestTable := []test_data.ProductTestTable{
		{
			Name:        "ERROR_UNKNOW_FIELD",
			QueryParams: "",
			Body:        `{"name":"Test-1", "description": "Description Test-1", "price": 2300, "category": "Middle", "lol": 2123}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.JsonUnknownFieldError.Code,
			ExpectedError:  internalErrors.JsonUnknownFieldError.Err,
		},

		{
			Name:        "ERROR_INVALID_FIELD_TYPE",
			QueryParams: "",
			Body:        `{"name":"Test-1", "description": "Description Test-1", "price": 2300, "category": 345}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.JsonTypeError.Code,
			ExpectedError:  internalErrors.JsonTypeError.Err,
		},

		{
			Name:           "ERROR_NEGATIVE_PRICE",
			Body:           `{"name":"Test-1", "description": "Description Test-1", "price": -2300, "category": "Middle"}`,
			Setup:          func() {},
			Except:         models.Product{},
			ExpectedStatus: internalErrors.NegativePriceError.Code,
			ExpectedError:  internalErrors.NegativePriceError.Err,
		},
	}

	test_data.StartTest(t, successTestTable, r, http.MethodPost)

	test_data.StartTest(t, errorsEmptyFieldTestTable, r, http.MethodPost)
	test_data.StartTest(t, errorsFieldTestTable, r, http.MethodPost)
}

func TestUpdate(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2026-01-01 12:00:00")
	expectedProduct := models.Product{
		Id:          1,
		Name:        "Update Test",
		Description: "Update Description",
		Price:       9000,
		Category:    "Update Category",
		CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockProductService(ctrl)
	h := NewProductHandler(service)

	r := chi.NewRouter()
	r.Put("/products/{id}", h.UpdateProduct)

	successTestTable := []test_data.ProductTestTable{
		{
			Name:        "OK",
			QueryParams: "1",
			Body:        `{"name":"Update Test", "description": "Update Description", "price": 9000, "category": "Update Category"}`,
			Setup: func() {
				service.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any(), 1).
					Return(expectedProduct, nil)
			},

			Except:         expectedProduct,
			ExpectedStatus: http.StatusOK,
			ExpectedError:  nil,
		},
	}

	errorsEmptyFieldTestTable := []test_data.ProductTestTable{
		{
			Name:        "ERROR_EMPTY_NAME",
			QueryParams: "1",
			Body:        `{"name":"", "description": "Update Description", "price": 9000, "category": "Update Category"}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.EmptyFieldsError.Code,
			ExpectedError:  internalErrors.EmptyFieldsError.Err,
		},

		{
			Name:        "ERROR_EMPTY_TEXT_FIELDS",
			QueryParams: "1",
			Body:        `{"name":"", "description": "", "price": 9000, "category": ""}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.EmptyFieldsError.Code,
			ExpectedError:  internalErrors.EmptyFieldsError.Err,
		},

		{
			Name:        "ERROR_EMPTY_BODY",
			QueryParams: "1",
			Body:        `{}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.EmptyFieldsError.Code,
			ExpectedError:  internalErrors.EmptyFieldsError.Err,
		},
	}
	errorsFieldTestTable := []test_data.ProductTestTable{
		{
			Name:        "ERROR_UNKNOW_FIELD",
			QueryParams: "1",
			Body:        `{"name":"Update Test", "description": "Update Description", "price": 9000, "category": "Update Category", "lol":23123}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.JsonUnknownFieldError.Code,
			ExpectedError:  internalErrors.JsonUnknownFieldError.Err,
		},

		{
			Name:        "ERROR_INVALID_FIELD_TYPE",
			QueryParams: "1",
			Body:        `{"name":"Update Test", "description": "Update Description", "price": 9000, "category": 2000}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.JsonTypeError.Code,
			ExpectedError:  internalErrors.JsonTypeError.Err,
		},

		{
			Name:        "ERROR_NEGATIVE_PRICE",
			QueryParams: "1",
			Body:        `{"name":"Update Test", "description": "Update Description", "price": -9000, "category": "Update Category"}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.NegativePriceError.Code,
			ExpectedError:  internalErrors.NegativePriceError.Err,
		},
	}
	errorsIdTestTable := []test_data.ProductTestTable{
		{
			Name:        "ERROR_NOT_FOUND_ID",
			QueryParams: "200",
			Body:        `{"name":"Update Test", "description": "Update Description", "price": 9000, "category": "Update Category"}`,
			Setup: func() {
				service.EXPECT().
					UpdateProduct(gomock.Any(), gomock.Any(), 200).
					Return(models.Product{}, pgx.ErrNoRows)
			},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.DataNotFound.Code,
			ExpectedError:  internalErrors.DataNotFound.Err,
		},

		{
			Name:        "ERROR_RANGE_ID",
			QueryParams: "200000000000000000000000000000000000000000000000",
			Body:        `{"name":"Update Test", "description": "Update Description", "price": 9000, "category": "Update Category"}`,
			Setup:       func() {},

			Except:         expectedProduct,
			ExpectedStatus: internalErrors.IdRangeError.Code,
			ExpectedError:  internalErrors.IdRangeError.Err,
		},
	}

	test_data.StartTest(t, successTestTable, r, http.MethodPut)

	test_data.StartTest(t, errorsEmptyFieldTestTable, r, http.MethodPut)
	test_data.StartTest(t, errorsFieldTestTable, r, http.MethodPut)
	test_data.StartTest(t, errorsIdTestTable, r, http.MethodPut)
}
