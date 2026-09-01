package service

import (
	"JuniarProject/internal/models"
	repository "JuniarProject/internal/repository/product"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProductService interface {
	CreatProduct(ctx context.Context, product models.Product) (models.Product, error)
	GetProducts(ctx context.Context) ([]models.Product, error)
}

type ProductServiceImpl struct {
	repo repository.ProductRepositoryImpl
}

func NewProductService(repo repository.ProductRepositoryImpl) ProductService {
	return &ProductServiceImpl{repo}
}

func (s *ProductServiceImpl) CreatProduct(ctx context.Context, product models.Product) (models.Product, error) {
	product.CreatedAt = pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductServiceImpl) GetProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetProducts(ctx)
}
