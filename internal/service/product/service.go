package service

import (
	"JuniarProject/internal/models"
	repository "JuniarProject/internal/repository/product"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type ProductService interface {
	CreateProduct(ctx context.Context, product models.Product) (models.Product, error)
	GetProducts(ctx context.Context, parameters models.ProductsGetQueryParameters) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.ProductRequest, id int) (models.Product, error)
	DeleteProduct(ctx context.Context, id int) (models.Product, error)
}

type ProductServiceImpl struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{repo}
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	product.CreatedAt = pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductServiceImpl) GetProducts(ctx context.Context, parameters models.ProductsGetQueryParameters) ([]models.Product, error) {
	return s.repo.GetProducts(ctx, parameters)
}

func (s *ProductServiceImpl) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, product models.ProductRequest, id int) (models.Product, error) {
	return s.repo.UpdateProduct(ctx, product, id)
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id int) (models.Product, error) {
	return s.repo.DeleteProduct(ctx, id)
}
