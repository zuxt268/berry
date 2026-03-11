package port

import (
	"context"
)

// OAuthUser OAuth認証で取得したユーザー情報
type OAuthUser struct {
	Sub     string
	Email   string
	Name    string
	Picture string
}

// OAuthAdapter OAuth操作のインターフェース
type OAuthAdapter interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*OAuthUser, error)
}