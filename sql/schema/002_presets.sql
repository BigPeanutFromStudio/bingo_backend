-- +goose Up

CREATE TABLE presets (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  events JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  owner_id VARCHAR(21) NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE presets;   