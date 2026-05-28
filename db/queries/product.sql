-- name: GetAllProducts :many
SELECT id, name, price, stock FROM products;

-- name: GetProductByID :one
SELECT id, name, price, stock FROM products
WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, price, stock)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $1, price = $2, stock = $3
WHERE id = $4
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
