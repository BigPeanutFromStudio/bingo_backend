-- name: CreateGame :one
INSERT INTO games (id, name, end_time, created_at, updated_at,  preset, admin_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetGame :one
SELECT * FROM games WHERE id = $1;