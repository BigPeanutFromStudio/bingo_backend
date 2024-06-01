-- name: CreateUser :one
INSERT INTO users (id, nickname, public_id, email, picture_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPublicID :one
SELECT * FROM users WHERE public_id = $1;

-- name: GetAllUsers :many
SELECT nickname, public_id, picture_url FROM users;

-- name: UpdateUser :one
UPDATE users SET nickname = $1 
WHERE id = $2
RETURNING *;

