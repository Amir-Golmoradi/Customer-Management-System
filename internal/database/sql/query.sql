-- name: CreateCustomer :one
INSERT INTO customers (
    name,
    email,
    password
)
VALUES ($1, $2, $3)
RETURNING
    id,
    name,
    email,
    password,
    created_at,
    updated_at;



-- name: GetCustomerByID :one
SELECT
    id,
    name,
    email,
    password,
    created_at,
    updated_at
FROM customers
WHERE id = $1
LIMIT 1;



-- name: GetCustomerByEmail :one
SELECT
    id,
    name,
    email,
    password,
    created_at,
    updated_at
FROM customers
WHERE email = $1
LIMIT 1;



-- name: ListCustomers :many
SELECT
    id,
    name,
    email,
    password,
    created_at,
    updated_at
FROM customers
ORDER BY id;



-- name: UpdateCustomer :one
UPDATE customers
SET
    id = $1,
    name = $2,
    email = $3,
    password = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    name,
    email,
    password,
    created_at,
    updated_at;



-- name: DeleteCustomerByEmail :execrows
DELETE FROM customers
WHERE email = $1;