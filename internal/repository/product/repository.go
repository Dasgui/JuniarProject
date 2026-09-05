package repository

import (
	"JuniarProject/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type ProductRepository interface {
	CreateProduct(ctx context.Context, product models.Product) (models.Product, error)
	GetProducts(ctx context.Context, parameters models.ProductsGetQueryParameters) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.ProductRequest, id int) (models.Product, error)
	DeleteProduct(ctx context.Context, id int) (models.Product, error)
}

type ProductRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (rep *ProductRepositoryImpl) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	var result models.Product

	row := rep.db.QueryRow(ctx,
		"INSERT INTO products (name, description, price, category, created_at) VALUES($1, $2, $3, $4, $5) RETURNING *",
		product.Name, product.Description, product.Price, product.Category, product.CreatedAt)

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)
	if err != nil {
		return result, err
	}

	return result, err
}

func (rep *ProductRepositoryImpl) GetProducts(ctx context.Context, parameters models.ProductsGetQueryParameters) ([]models.Product, error) {

	var query = "SELECT * FROM products WHERE 1=1"
	args := make([]any, 0)

	var rows pgx.Rows
	var err error

	var result []models.Product

	query, args = addParameters(query, args, parameters)

	rows, err = rep.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Category, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, err
}

func addParameters(query string, args []any, parameters models.ProductsGetQueryParameters) (string, []any) {
	argCount := 1

	if parameters.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, parameters.Category)
		argCount++
	}

	if parameters.PriceFrom != 0 {
		query += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, parameters.PriceFrom)
		argCount++
	}

	if parameters.PriceTo != 0 {
		query += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, parameters.PriceTo)
		argCount++
	}

	if parameters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, parameters.Limit)
		argCount++
	}

	if parameters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, parameters.Offset)
		argCount++
	}

	return query, args
}

func (rep *ProductRepositoryImpl) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	var result models.Product

	row := rep.db.QueryRow(ctx, "SELECT * FROM products WHERE id = $1", id)

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)

	if err != nil {
		return result, err
	}
	return result, err
}

func (rep *ProductRepositoryImpl) UpdateProduct(ctx context.Context, product models.ProductRequest, id int) (models.Product, error) {
	var result models.Product

	row := rep.db.QueryRow(ctx, "UPDATE products SET name = $1, description = $2, price = $3, category = $4 WHERE id = $5 RETURNING *",
		product.Name, product.Description, product.Price, product.Category, id)

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)
	if err != nil {
		return result, err
	}

	return result, err
}

func (rep *ProductRepositoryImpl) DeleteProduct(ctx context.Context, id int) (models.Product, error) {
	var result models.Product

	row := rep.db.QueryRow(ctx, "DELETE FROM products WHERE id = $1 RETURNING *", id)

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)

	if err != nil {
		return result, err
	}

	return result, err
}
