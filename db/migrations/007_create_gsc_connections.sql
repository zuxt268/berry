-- +migrate Up
CREATE TABLE gsc_connections (
    id BIGINT NOT NULL AUTO_INCREMENT,
    uid VARCHAR(36) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    site_url VARCHAR(512) NOT NULL COMMENT 'Search Consoleのサイト URL',
    refresh_token TEXT NOT NULL,
    connected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    disconnected_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_uid (uid),
    INDEX idx_user_id (user_id),
    CONSTRAINT fk_gsc_connections_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS gsc_connections;