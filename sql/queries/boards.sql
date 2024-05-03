-- name: CreatePreset :one
INSERT INTO presets (id, name, events, created_at, updated_at, owner_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserPresets :many
SELECT * FROM presets WHERE owner_id = $1;