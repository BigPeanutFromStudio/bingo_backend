-- name: CreateUser :one
INSERT INTO users (id, nickname, email, picture_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET nickname = $1 
WHERE id = $2
RETURNING *;