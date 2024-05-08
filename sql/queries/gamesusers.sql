-- name: CreateGameUser :one
INSERT INTO games_users (id, user_id, game_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetGameUsers :many
SELECT * FROM games_users WHERE user_id = $1;

-- name: DeleteGameUsers :exec
DELETE FROM games_users WHERE id = $1 AND user_id = $2; 