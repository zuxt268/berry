-- +migrate Up
CREATE TABLE operator_sessions (
    id BIGSERIAL PRIMARY KEY,
    uid VARCHAR(36) NOT NULL UNIQUE,
    operator_id BIGINT NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS operator_sessions;