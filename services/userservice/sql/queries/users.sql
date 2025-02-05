-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, phone, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: SetUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;
