-- +goose Up
CREATE TABLE IF NOT EXISTS Products (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price float NOT NULL,
    category TEXT NOT NULL,
    created_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS Products;
