-- +goose Up
CREATE TABLE IF NOT EXISTS Products (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price float NOT NULL,
    category TEXT NOT NULL,
    created_at timestamp NOT NULL
);
INSERT INTO Products (name, description, price, category, created_at)
VALUES('Test-1', 'Description Test-1', 2300, 'Test', '2026-09-02T11:37:16.77623'),
('Test-2', 'Description Test-2', 6500, 'Test', '2026-09-03T11:37:16.77623'),
('Test-3', 'Description Test-3', 8300, 'Test-3', '2026-09-04T11:37:16.77623');

-- +goose Down
DROP TABLE IF EXISTS Products;
