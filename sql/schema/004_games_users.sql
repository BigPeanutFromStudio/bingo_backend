-- +goose Up

CREATE TABLE games_users (
  id UUID PRIMARY KEY,
  game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  user_id VARCHAR(21) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE (game_id, user_id)
);

-- +goose Down

DROP TABLE games_users;   