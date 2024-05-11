-- name: CreateUser :one
INSERT INTO users (id, nickname, email ,refresh_token, picture_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;