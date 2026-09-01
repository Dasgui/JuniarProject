package repository

import (
	"JuniarProject/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	CreatProduct(ctx context.Context, product models.Product) (models.Product, error)
	GetProducts(ctx context.Context) ([]models.Product, error)
}

type ProductRepositoryImpl struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) ProductRepositoryImpl {
	return ProductRepositoryImpl{db: db}
}

func (rep *ProductRepositoryImpl) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	row := rep.db.QueryRow(ctx,
		"INSERT INTO products (name, description, price, category, created_at) VALUES($1, $2, $3, $4, $5) RETURNING *",
		product.Name, product.Description, product.Price, product.Category, product.CreatedAt)

	var result models.Product

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)
	if err != nil {
		return result, err
	}

	return result, err
}

func (rep *ProductRepositoryImpl) GetProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := rep.db.Query(ctx, "SELECT * FROM products")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Product

	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Category, &product.CreatedAt)
		result = append(result, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, err
}
