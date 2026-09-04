package handler

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/lib/test_data"
	"JuniarProject/internal/models"
	mock_service "JuniarProject/internal/service/product/mocks"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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

	successTestTable := []test_data.TestTable{
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

	errorTestTable := []test_data.TestTable{
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
	test_data.StartTest(t, successTestTable, h.GetProducts)
	test_data.StartTest(t, errorTestTable, h.GetProducts)
}
