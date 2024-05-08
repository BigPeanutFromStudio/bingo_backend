-- +goose Up

CREATE TABLE games (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  end_time TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  preset UUID NOT NULL REFERENCES presets(id) ON DELETE CASCADE,
  admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE games;   