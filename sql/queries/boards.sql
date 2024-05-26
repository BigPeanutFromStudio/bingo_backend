-- name: CreateBoard :one
INSERT INTO boards (id, events, created_at, updated_at, game_id, owner_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserBoards :many
SELECT * FROM boards WHERE owner_id = $1;