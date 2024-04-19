-- +goose UP

CREATE TABLE users (
  id UUID PRIMARY KEY,
  nickname TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose DOWN

DROP TABLE users;   