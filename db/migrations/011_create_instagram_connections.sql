-- +migrate Up
CREATE TABLE instagram_connections (
    id BIGINT NOT NULL AUTO_INCREMENT,
    uid VARCHAR(36) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    instagram_business_account_id VARCHAR(255) NOT NULL COMMENT 'IGビジネスアカウントID',
    facebook_page_id VARCHAR(255) NOT NULL COMMENT '紐づくFacebookページID',
    access_token TEXT NOT NULL COMMENT '長期アクセストークン',
    token_expires_at TIMESTAMP NULL DEFAULT NULL COMMENT 'トークン有効期限',
    connected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    disconnected_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_uid (uid),
    INDEX idx_user_id (user_id),
    INDEX idx_token_expires (token_expires_at),
    CONSTRAINT fk_instagram_connections_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS instagram_connections;