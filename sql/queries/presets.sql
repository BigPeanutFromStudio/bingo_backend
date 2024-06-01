-- name: CreatePreset :one
INSERT INTO presets (id, name, events, created_at, updated_at, owner_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserPresets :many
SELECT * FROM presets WHERE owner_id = $1;

-- name: GetUserPresetByID :one
SELECT * FROM presets WHERE id = $1;

-- name: UpdatePresetEvents :one
UPDATE presets SET events = $1
WHERE id = $2
RETURNING *;