-- +migrate Up
CREATE TABLE gbp_daily_reports (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    gbp_connection_id BIGINT NOT NULL,
    report_date DATE NOT NULL COMMENT '集計対象日',
    profile_views INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'ビジネスプロフィール閲覧数',
    phone_calls INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '電話タップ数',
    direction_requests INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'ルート検索数',
    photo_views INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '写真閲覧数',
    review_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'クチコミ数',
    average_rating DECIMAL(3,2) NOT NULL DEFAULT 0 COMMENT '平均評価',
    search_query_breakdown JSON COMMENT '検索クエリ別内訳',
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'バッチ取得時刻',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_connection_date (gbp_connection_id, report_date),
    INDEX idx_report_date (report_date),
    CONSTRAINT fk_gbp_daily_reports_connection FOREIGN KEY (gbp_connection_id) REFERENCES gbp_connections(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS gbp_daily_reports;