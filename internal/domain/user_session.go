package domain

import "time"

type UserSession struct {
	ID           uint64
	UID          string
	UserID       uint64
	SessionToken string
	IPAddress    string
	UserAgent    string
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// IsExpired セッションが期限切れかチェック
func (s *UserSession) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

// Refresh セッションを更新
func (s *UserSession) Refresh(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
	s.UpdatedAt = time.Now()
}
