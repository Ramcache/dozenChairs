-- +goose Up
CREATE TABLE user_sessions (
                               id UUID PRIMARY KEY,
                               user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                               token_hash TEXT NOT NULL,
                               user_agent TEXT,
                               ip_address TEXT,
                               expires_at TIMESTAMP NOT NULL,
                               created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);

-- +goose Down
DROP TABLE user_sessions;