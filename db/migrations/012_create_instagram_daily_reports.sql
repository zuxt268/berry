-- +migrate Up
CREATE TABLE instagram_daily_reports (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    instagram_connection_id BIGINT NOT NULL,
    report_date DATE NOT NULL COMMENT '集計対象日',
    follower_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'フォロワー数',
    impressions INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'インプレッション数',
    reach INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'リーチ数',
    profile_views INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'プロフィール閲覧数',
    website_clicks INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'ウェブサイトクリック数',
    post_engagement JSON COMMENT '投稿別エンゲージメント内訳',
    audience_demographics JSON COMMENT 'フォロワー属性(性別年齢・都市・国)',
    stories_insights JSON COMMENT 'ストーリーズインサイト',
    fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'バッチ取得時刻',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_connection_date (instagram_connection_id, report_date),
    INDEX idx_report_date (report_date),
    CONSTRAINT fk_instagram_daily_reports_connection FOREIGN KEY (instagram_connection_id) REFERENCES instagram_connections(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS instagram_daily_reports;