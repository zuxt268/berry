-- +migrate Up
CREATE TABLE ga4_daily_reports (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    ga4_connection_id BIGINT NOT NULL,
    report_date DATE NOT NULL COMMENT '集計対象日',
    sessions INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '総セッション数',
    total_users INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'ユニークユーザー数',
    bounce_rate DECIMAL(5,4) NOT NULL DEFAULT 0 COMMENT '直帰率',
    avg_session_duration DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '平均セッション時間(秒)',
    conversions INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'コンバージョン数',
    channel_breakdown JSON COMMENT '流入経路別内訳',
    device_breakdown JSON COMMENT 'デバイス別アクセス',
    page_breakdown JSON COMMENT 'ページ別PV',
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'バッチ取得時刻',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_connection_date (ga4_connection_id, report_date),
    INDEX idx_report_date (report_date),
    CONSTRAINT fk_ga4_daily_reports_connection FOREIGN KEY (ga4_connection_id) REFERENCES ga4_connections(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS ga4_daily_reports;