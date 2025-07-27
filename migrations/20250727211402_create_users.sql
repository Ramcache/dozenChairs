-- +goose Up
CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       username TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       role VARCHAR DEFAULT 'user',
                       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;