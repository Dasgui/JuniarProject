package repository

import (
	"JuniarProject/internal/internalErrors"
	"JuniarProject/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	CreatProduct(ctx context.Context, product models.Product) (models.Product, error)
	GetProducts(ctx context.Context, category string, priceFrom float64, priceTo float64, limit int, offset int) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product, id int) (models.Product, error)
	DeleteProduct(ctx context.Context, id int) (models.Product, error)
}

type ProductRepositoryImpl struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) ProductRepositoryImpl {
	return ProductRepositoryImpl{db: db}
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

func (rep *ProductRepositoryImpl) GetProducts(ctx context.Context, category string, priceFrom float64, priceTo float64,
	limit int, offset int) ([]models.Product, error) {

	var query = "SELECT * FROM products WHERE 1=1"
	args := make([]any, 0)
	argCount := 1

	var rows pgx.Rows
	var err error

	var result []models.Product

	if category != "" {
		query = query + fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, category)
		argCount++
	}

	if priceFrom != 0 {
		query = query + fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, priceFrom)
		argCount++
	}

	if priceTo != 0 {
		query = query + fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, priceTo)
		argCount++
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
		argCount++
	}

	rows, err = rep.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (rep *ProductRepositoryImpl) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	var result models.Product

	row := rep.db.QueryRow(ctx, "SELECT * FROM products WHERE id = $1", id)

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)

	if err != nil {
		return result, internalErrors.HandleDbError(err)
	}
	return result, err
}

func (rep *ProductRepositoryImpl) UpdateProduct(ctx context.Context, product models.Product, id int) (models.Product, error) {
	var result models.Product

	row := rep.db.QueryRow(ctx, "UPDATE products SET name = $1, description = $2, price = $3, category = $4, created_at = $5 WHERE id = $6 RETURNING *",
		product.Name, product.Description, product.Price, product.Category, product.CreatedAt, id)

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)
	if err != nil {
		return result, internalErrors.HandleDbError(err)
	}

	return result, err
}

func (rep *ProductRepositoryImpl) DeleteProduct(ctx context.Context, id int) (models.Product, error) {
	var result models.Product

	row := rep.db.QueryRow(ctx, "DELETE FROM products WHERE id = $1 RETURNING *", id)

	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Price, &result.Category, &result.CreatedAt)

	if err != nil {
		return result, internalErrors.HandleDbError(err)
	}

	return result, err
}
