package service

//
//import (
//	"JuniarProject/internal/models"
//	mock_repository "JuniarProject/internal/repository/product/mocks"
//	"context"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"github.com/jackc/pgx/v5/pgtype"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
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
//			Category:    "Test",
//			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
//		},
//		{
//			Id:          2,
//			Name:        "Test-2",
//			Description: "Description Test-2",
//			Price:       6500,
//			Category:    "Test",
//			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
//		},
//		{
//			Id:          1,
//			Name:        "Test-1",
//			Description: "Description Test-1",
//			Price:       8300,
//			Category:    "Test",
//			CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
//		},
//	}
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	repository := mock_repository.NewMockProductRepository(ctrl)
//	service := NewProductService(repository)
//	ctx := context.Background()
//
//	testTable := []struct {
//		name string
//
//		category  string
//		priceFrom float64
//		priceTo   float64
//		limit     int
//		offset    int
//
//		setup   func()
//		except  []models.Product
//		wantErr bool
//	}{
//		{
//			name:      "OK",
//			category:  "",
//			priceFrom: 0,
//			priceTo:   0,
//			limit:     0,
//			offset:    0,
//			setup: func() {
//				repository.EXPECT().
//					GetProducts(ctx, "", 0.0, 0.0, 0, 0).
//					Return(expectedProducts, nil)
//			},
//			except:  expectedProducts,
//			wantErr: false,
//		},
//		{
//			name:      "NOT_FIND_CATEGORY",
//			category:  "Ultra Category",
//			priceFrom: 0,
//			priceTo:   0,
//			limit:     0,
//			offset:    0,
//			setup: func() {
//				repository.EXPECT().
//					GetProducts(ctx, "Ultra Category", 0.0, 0.0, 0, 0).
//					Return([]models.Product{}, nil)
//			},
//			except:  []models.Product{},
//			wantErr: false,
//		},
//	}
//
//	for _, tt := range testTable {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setup()
//
//			got, err := service.GetProducts(ctx, tt.category, tt.priceFrom, tt.priceTo, tt.limit, tt.offset)
//
//			if tt.wantErr {
//				assert.Error(t, err)
//				return
//			}
//
//			require.NoError(t, err)
//			assert.Equal(t, tt.except, got)
//		})
//	}
//
//}
