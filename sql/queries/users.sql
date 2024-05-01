-- name: CreateUser :one
INSERT INTO users (id, nickname, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByToken :one
SELECT * FROM users WHERE token = $1;