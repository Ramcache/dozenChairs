-- +goose Up
CREATE TABLE IF NOT EXISTS chairs (
                                      id VARCHAR(32) PRIMARY KEY,
                                      type VARCHAR(16) NOT NULL,
                                      category VARCHAR(16) NOT NULL,
                                      title VARCHAR(128) NOT NULL,
                                      slug VARCHAR(128) NOT NULL UNIQUE,
                                      description TEXT,
                                      price INT NOT NULL,
                                      old_price INT,
                                      in_stock BOOLEAN NOT NULL,
                                      unit_count INT,
                                      images JSONB NOT NULL,
                                      attributes JSONB NOT NULL,
                                      tags JSONB NOT NULL,
                                      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS chairs;