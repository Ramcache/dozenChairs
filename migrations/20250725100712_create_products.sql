-- +goose Up
CREATE TABLE products (
                          id TEXT PRIMARY KEY,
                          type TEXT NOT NULL,
                          category TEXT NOT NULL,
                          title TEXT NOT NULL,
                          slug TEXT NOT NULL,
                          description TEXT,
                          price INTEGER NOT NULL,
                          old_price INTEGER,
                          in_stock BOOLEAN NOT NULL,
                          unit_count INTEGER,
                          images JSONB,
                          attributes JSONB,
                          includes JSONB,
                          tags JSONB,
                          created_at TIMESTAMP NOT NULL,
                          updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE products;