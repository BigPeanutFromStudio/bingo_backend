// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: gamesusers.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createGameUser = `-- name: CreateGameUser :one
INSERT INTO games_users (id, user_id, game_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, game_id, user_id, created_at, updated_at
`

type CreateGameUserParams struct {
	ID        uuid.UUID
	UserID    string
	GameID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateGameUser(ctx context.Context, arg CreateGameUserParams) (GamesUser, error) {
	row := q.db.QueryRowContext(ctx, createGameUser,
		arg.ID,
		arg.UserID,
		arg.GameID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i GamesUser
	err := row.Scan(
		&i.ID,
		&i.GameID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteGameUsers = `-- name: DeleteGameUsers :exec
DELETE FROM games_users WHERE id = $1 AND user_id = $2
`

type DeleteGameUsersParams struct {
	ID     uuid.UUID
	UserID string
}

func (q *Queries) DeleteGameUsers(ctx context.Context, arg DeleteGameUsersParams) error {
	_, err := q.db.ExecContext(ctx, deleteGameUsers, arg.ID, arg.UserID)
	return err
}

const getGameUsers = `-- name: GetGameUsers :many
SELECT id, game_id, user_id, created_at, updated_at FROM games_users WHERE user_id = $1
`

func (q *Queries) GetGameUsers(ctx context.Context, userID string) ([]GamesUser, error) {
	rows, err := q.db.QueryContext(ctx, getGameUsers, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GamesUser
	for rows.Next() {
		var i GamesUser
		if err := rows.Scan(
			&i.ID,
			&i.GameID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}