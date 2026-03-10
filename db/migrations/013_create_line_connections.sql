-- +migrate Up
CREATE TABLE line_connections (
    id BIGINT NOT NULL AUTO_INCREMENT,
    uid VARCHAR(36) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    channel_id VARCHAR(255) NOT NULL COMMENT 'LINE Messaging APIチャンネルID',
    channel_secret VARCHAR(255) NOT NULL COMMENT 'チャンネルシークレット',
    channel_access_token TEXT NOT NULL COMMENT 'チャンネルアクセストークン',
    channel_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '管理用表示名',
    bot_user_id VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'BotのユーザーID',
    connected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    disconnected_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_uid (uid),
    INDEX idx_user_id (user_id),
    UNIQUE KEY uk_user_channel (user_id, channel_id),
    CONSTRAINT fk_line_connections_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS line_connections;