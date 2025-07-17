-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(64) NOT NULL UNIQUE,
                                     password_hash VARCHAR(255) NOT NULL,
                                     email VARCHAR(255) NOT NULL UNIQUE,
                                     full_name VARCHAR(255),
                                     phone VARCHAR(32),
                                     address VARCHAR(255),
                                     role VARCHAR(32) NOT NULL DEFAULT 'user',
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     email_verified BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE IF EXISTS users;