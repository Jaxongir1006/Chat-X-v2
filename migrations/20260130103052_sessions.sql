-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),

    access_token VARCHAR(500) NOT NULL UNIQUE,
    access_token_expires_at TIMESTAMPTZ NOT NULL,

    refresh_token VARCHAR(500) NOT NULL UNIQUE,
    refresh_token_expires_at TIMESTAMPTZ NOT NULL,

    ip_address VARCHAR(50) NOT NULL,
    user_agent VARCHAR(500) NOT NULL,
    device VARCHAR(250) NOT NULL,

    revoked_at TIMESTAMPTZ,

    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_refresh_exp
    ON sessions (user_id, refresh_token_expires_at);

CREATE INDEX idx_sessions_refresh_exp
    ON sessions (refresh_token_expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE sessions;
-- +goose StatementEnd
