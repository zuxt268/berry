package adapter

import (
	"context"
	"fmt"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GBPOAuthResult GBP OAuth認証で取得したトークン情報
type GBPOAuthResult struct {
	RefreshToken string
	AccessToken  string
}

// GBPOAuthAdapter GBP OAuth操作のインターフェース
//
//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package mock
type GBPOAuthAdapter interface {
	// GetAuthURL stateパラメータ付きのOAuth認証URLを返す
	GetAuthURL(state string) string
	// ExchangeCode 認証コードをトークン情報と交換
	ExchangeCode(ctx context.Context, code string) (*GBPOAuthResult, error)
}

type gbpOAuthClient struct {
	config *oauth2.Config
}

func NewGBPOAuthClient(redirectURL string) GBPOAuthAdapter {
	oauthConfig := &oauth2.Config{
		ClientID:     config.Env.GoogleClientID,
		ClientSecret: config.Env.GoogleClientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/business.manage",
		},
		Endpoint: google.Endpoint,
	}
	return &gbpOAuthClient{config: oauthConfig}
}

// GetAuthURL stateパラメータ付きのOAuth認証URLを返す（refresh_token取得のためprompt=consentを付与）
func (c *gbpOAuthClient) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
}

// ExchangeCode 認証コードをトークン情報と交換
func (c *gbpOAuthClient) ExchangeCode(ctx context.Context, code string) (*GBPOAuthResult, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrOAuthTokenExchange, err)
	}

	if token.RefreshToken == "" {
		return nil, fmt.Errorf("%w: refresh token not provided", domain.ErrOAuthTokenExchange)
	}

	return &GBPOAuthResult{
		RefreshToken: token.RefreshToken,
		AccessToken:  token.AccessToken,
	}, nil
}
