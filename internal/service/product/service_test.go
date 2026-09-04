package service

//
//import (
//	"JuniarProject/internal/lib/test_data"
//	"JuniarProject/internal/models"
//	mock_repository "JuniarProject/internal/repository/product/mocks"
//	"net/http"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"github.com/jackc/pgx/v5/pgtype"
//)
//
//func TestGet(t *testing.T) {
//	createdAt, _ := time.Parse(time.RFC3339, "2026-01-01 12:00:00")
//	expectedProducts := []models.Product{
//		{
//			Id:          1,
//			Name:        "Test-1",
//			Description: "Description Test-1",
//			Price:       2300,
//			Category:    "Middle",
//			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
//		},
//		{
//			Id:          2,
//			Name:        "Test-2",
//			Description: "Description Test-2",
//			Price:       6500,
//			Category:    "Electro",
//			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
//		},
//		{
//			Id:          3,
//			Name:        "Test-1",
//			Description: "Description Test-1",
//			Price:       8300,
//			Category:    "Electro",
//			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
//		},
//	}
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	repository := mock_repository.NewMockProductRepository(ctrl)
//
//	successTestTable := []test_data.TestServiceTable{
//		{
//			Name:       "OK_NO_FILTER",
//			Parameters: models.ProductsGetQueryParameters{},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{}).
//					Return(expectedProducts, nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  len(expectedProducts),
//		},
//
//		{
//			Name:       "OK_FILTER_CATEGORY",
//			Parameters: models.ProductsGetQueryParameters{Category: "Electro"},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						Category: "Electro",
//					}).
//					Return(expectedProducts[:2], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  2,
//		},
//
//		{
//			Name: "OK_FILTER_PRICE_FROM",
//			Parameters: models.ProductsGetQueryParameters{
//				PriceFrom: 6000.0,
//			},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						PriceFrom: 6000.0,
//					}).
//					Return(expectedProducts[1:], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  2,
//		},
//
//		{
//			Name: "OK_FILTER_PRICE_TO",
//			Parameters: models.ProductsGetQueryParameters{
//				PriceTo: 7000.0,
//			},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						PriceTo: 7000.0,
//					}).
//					Return(expectedProducts[:2], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  2,
//		},
//
//		{
//			Name: "OK_FILTER_PRICE_FROM_TO",
//			Parameters: models.ProductsGetQueryParameters{
//				PriceFrom: 6000.0, PriceTo: 7000.0,
//			},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						PriceFrom: 6000.0, PriceTo: 7000.0,
//					}).
//					Return(expectedProducts[2:3], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  1,
//		},
//
//		{
//			Name: "OK_FILTER_CATEGORY_PRICE_FROM_TO",
//			Parameters: models.ProductsGetQueryParameters{
//				Category: "Middle", PriceFrom: 2000.0, PriceTo: 10000.0,
//			},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						Category: "Middle", PriceFrom: 2000.0, PriceTo: 10000.0,
//					}).
//					Return(expectedProducts[1:2], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  1,
//		},
//
//		{
//			Name: "OK_FILTER_LIMIT",
//			Parameters: models.ProductsGetQueryParameters{
//				Limit: 1,
//			},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						Limit: 1,
//					}).
//					Return(expectedProducts[1:2], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  1,
//		},
//
//		{
//			Name: "OK_FILTER_OFFSET",
//			Parameters: models.ProductsGetQueryParameters{
//				Offset: 2,
//			},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						Offset: 2,
//					}).
//					Return(expectedProducts[:2], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  2,
//		},
//
//		{
//			Name: "OK_FILTER_LIMIT_OFFSET",
//			Parameters: models.ProductsGetQueryParameters{
//				Limit: 2, Offset: 1,
//			},
//			Setup: func() {
//				repository.EXPECT().
//					GetProducts(gomock.Any(), models.ProductsGetQueryParameters{
//						Limit: 2, Offset: 1,
//					}).
//					Return(expectedProducts[1:3], nil)
//			},
//
//			Except:         expectedProducts,
//			ExpectedStatus: http.StatusOK,
//			ExpectedError:  nil,
//			ExpectedCount:  2,
//		},
//	}
//
//	//errorTestTable := []test_data.TestServiceTable{
//	//	{
//	//		Name:       "ERROR_INVALID_PARAMETRES",
//	//		Parameters: models.ProductsGetQueryParameters{},
//	//		Setup: func() {
//	//			repository.EXPECT().
//	//				GetProducts(gomock.Any(), models.ProductsGetQueryParameters{}).
//	//				Return(nil, assert.AnError)
//	//		},
//	//		Except:        nil,
//	//		ExpectedError: assert.AnError,
//	//		ExpectedCount: 0,
//	//	},
//	//}
//
//	test_data.StartRepositoryTest(t, successTestTable, repository.GetProducts)
//}
