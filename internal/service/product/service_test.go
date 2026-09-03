package service

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/models"
	mock_repository "JuniarProject/internal/repository/product/mocks"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2026-01-01 12:00:00")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_repository.NewMockProductRepository(ctrl)
	service := NewProductService(repository)
	ctx := context.Background()

	expectedProduct := models.Product{
		Id:          3,
		Name:        "Test-3",
		Description: "Description Test-3",
		Price:       8300,
		Category:    "Test-3",
		CreatedAt:   pgtype.Timestamp{Time: createdAt, Valid: true},
	}

	testTable := []struct {
		name    string
		id      int
		setup   func()
		except  models.Product
		wantErr bool
	}{
		{
			name: "OK",
			id:   3,
			setup: func() {
				repository.EXPECT().
					GetProductByID(ctx, 3).
					Return(expectedProduct, nil)
			},
			except:  expectedProduct,
			wantErr: false,
		},
		{
			name: "NOT_FOUND_ERROR",
			id:   999,
			setup: func() {
				repository.EXPECT().
					GetProductByID(ctx, 999).
					Return(models.Product{}, internalErrors.DataNotFound)
			},
			except:  models.Product{},
			wantErr: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := service.GetProductByID(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.except, got)
		})
	}

}
