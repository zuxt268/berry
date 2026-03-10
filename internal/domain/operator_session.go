package domain

import "time"

// OperatorSession オペレーターセッションエンティティ
type OperatorSession struct {
	ID           int64
	UID          string
	OperatorID   int64
	SessionToken string
	IPAddress    string
	UserAgent    string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

// IsExpired セッションが期限切れかチェック
func (s *OperatorSession) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
