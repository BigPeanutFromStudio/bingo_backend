// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: games.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createGame = `-- name: CreateGame :one
INSERT INTO games (id, name, end_time, created_at, updated_at,  preset, admin_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, end_time, created_at, updated_at, preset, admin_id
`

type CreateGameParams struct {
	ID        uuid.UUID
	Name      string
	EndTime   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Preset    uuid.UUID
	AdminID   string
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (Game, error) {
	row := q.db.QueryRowContext(ctx, createGame,
		arg.ID,
		arg.Name,
		arg.EndTime,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Preset,
		arg.AdminID,
	)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.EndTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Preset,
		&i.AdminID,
	)
	return i, err
}

const getAdminedGames = `-- name: GetAdminedGames :many
SELECT id, name, end_time, created_at, updated_at, preset, admin_id FROM games WHERE admin_id = $1
`

func (q *Queries) GetAdminedGames(ctx context.Context, adminID string) ([]Game, error) {
	rows, err := q.db.QueryContext(ctx, getAdminedGames, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Game
	for rows.Next() {
		var i Game
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.EndTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Preset,
			&i.AdminID,
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

const getGame = `-- name: GetGame :one
SELECT id, name, end_time, created_at, updated_at, preset, admin_id FROM games WHERE id = $1
`

func (q *Queries) GetGame(ctx context.Context, id uuid.UUID) (Game, error) {
	row := q.db.QueryRowContext(ctx, getGame, id)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.EndTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Preset,
		&i.AdminID,
	)
	return i, err
}
