-- +goose Up

CREATE TABLE presets (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  events JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  owner_id UUID NOT NULL REFERENCES users(id)
);

-- +goose Down

DROP TABLE presets;   