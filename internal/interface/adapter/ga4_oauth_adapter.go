package adapter

import (
	"context"
	"fmt"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type ga4OAuthClient struct {
	config *oauth2.Config
}

func NewGA4OAuthClient(redirectURL string) port.GA4OAuthAdapter {
	oauthConfig := &oauth2.Config{
		ClientID:     config.Env.GoogleClientID,
		ClientSecret: config.Env.GoogleClientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/analytics.readonly",
		},
		Endpoint: google.Endpoint,
	}
	return &ga4OAuthClient{config: oauthConfig}
}

// GetAuthURL stateパラメータ付きのOAuth認証URLを返す（refresh_token取得のためprompt=consentを付与）
func (c *ga4OAuthClient) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
}

// ExchangeCode 認証コードをトークン情報と交換
func (c *ga4OAuthClient) ExchangeCode(ctx context.Context, code string) (*port.GA4OAuthResult, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrOAuthTokenExchange, err)
	}

	if token.RefreshToken == "" {
		return nil, fmt.Errorf("%w: refresh token not provided", domain.ErrOAuthTokenExchange)
	}

	return &port.GA4OAuthResult{
		RefreshToken: token.RefreshToken,
		AccessToken:  token.AccessToken,
	}, nil
}
