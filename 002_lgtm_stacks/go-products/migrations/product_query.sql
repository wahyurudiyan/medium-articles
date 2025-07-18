-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetProductBySKU :one
SELECT * FROM products
WHERE sku = $1 AND deleted_at IS NULL;

-- name: ListProducts :many
SELECT * FROM products
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateProducts :copyfrom
INSERT INTO products (sku, name, quantity, price)
VALUES ($1, $2, $3, $4);

-- name: UpdateProduct :one
UPDATE products
SET
  name = $1,
  quantity = $2,
  price = $3,
  updated_at = timezone('utc', now())
WHERE sku = $4 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteProduct :one
UPDATE products
SET
  deleted_at = timezone('utc', now()),
  updated_at = timezone('utc', now())
WHERE sku = $1 AND deleted_at IS NULL
RETURNING id, sku;
