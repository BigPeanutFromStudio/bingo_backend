-- +goose Up

CREATE TABLE users (
  id VARCHAR(21) PRIMARY KEY,
  nickname TEXT NOT NULL,
  email TEXT NOT NULL,
  refresh_token TEXT NOT NULL,
  picture_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down

DROP TABLE users;   