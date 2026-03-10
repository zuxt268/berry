-- +migrate Up
CREATE TABLE line_daily_reports (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    line_connection_id BIGINT NOT NULL,
    report_date DATE NOT NULL COMMENT '集計対象日',
    followers INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '友だち数',
    targeted_reaches INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'ターゲットリーチ数',
    blocks INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'ブロック数',
    message_delivery JSON COMMENT 'メッセージ配信統計',
    demographic JSON COMMENT '友だち属性（性別・年齢・地域）',
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'バッチ取得時刻',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_connection_date (line_connection_id, report_date),
    INDEX idx_report_date (report_date),
    CONSTRAINT fk_line_daily_reports_connection FOREIGN KEY (line_connection_id) REFERENCES line_connections(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS line_daily_reports;