-- +migrate Up
CREATE TABLE gsc_daily_reports (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    gsc_connection_id BIGINT NOT NULL,
    report_date DATE NOT NULL COMMENT '集計対象日',
    impressions INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '検索表示回数',
    clicks INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'クリック数',
    ctr DECIMAL(7,6) NOT NULL DEFAULT 0 COMMENT 'クリック率',
    average_position DECIMAL(8,2) NOT NULL DEFAULT 0 COMMENT '平均掲載順位',
    keyword_breakdown JSON COMMENT 'キーワード別内訳',
    page_breakdown JSON COMMENT 'ページ別パフォーマンス',
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'バッチ取得時刻',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_connection_date (gsc_connection_id, report_date),
    INDEX idx_report_date (report_date),
    CONSTRAINT fk_gsc_daily_reports_connection FOREIGN KEY (gsc_connection_id) REFERENCES gsc_connections(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS gsc_daily_reports;