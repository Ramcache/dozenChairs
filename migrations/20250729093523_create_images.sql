-- +goose Up
CREATE TABLE images (
                        id UUID PRIMARY KEY,
                        product_id TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
                        url TEXT NOT NULL,
                        filename TEXT NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS images;
