package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type instagramOAuthClient struct {
	config *oauth2.Config
}

func NewInstagramOAuthClient(redirectURL string) port.InstagramOAuthAdapter {
	oauthConfig := &oauth2.Config{
		ClientID:     config.Env.MetaAppID,
		ClientSecret: config.Env.MetaAppSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"instagram_basic",
			"instagram_manage_insights",
			"pages_show_list",
		},
		Endpoint: facebook.Endpoint,
	}
	return &instagramOAuthClient{config: oauthConfig}
}

// GetAuthURL stateパラメータ付きのOAuth認証URLを返す
func (c *instagramOAuthClient) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state)
}

// ExchangeCode 認証コードを長期トークン+IGビジネスアカウント情報と交換
func (c *instagramOAuthClient) ExchangeCode(ctx context.Context, code string) (*port.InstagramOAuthResult, error) {
	// 1. 短期トークンを取得
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrOAuthTokenExchange, err)
	}

	// 2. 短期トークンを長期トークンに交換
	longLivedToken, expiresAt, err := c.exchangeLongLivedToken(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInstagramTokenExchange, err)
	}

	// 3. Facebookページ一覧からIGビジネスアカウントIDを取得
	pageID, igAccountID, err := c.fetchInstagramBusinessAccount(ctx, longLivedToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInstagramAPICall, err)
	}

	return &port.InstagramOAuthResult{
		AccessToken:                longLivedToken,
		TokenExpiresAt:             *expiresAt,
		InstagramBusinessAccountID: igAccountID,
		FacebookPageID:             pageID,
	}, nil
}

// RefreshLongLivedToken 長期トークンを更新
func (c *instagramOAuthClient) RefreshLongLivedToken(ctx context.Context, currentToken string) (string, *time.Time, error) {
	url := fmt.Sprintf(
		"https://graph.facebook.com/v21.0/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		c.config.ClientID, c.config.ClientSecret, currentToken,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", domain.ErrInstagramTokenRefresh, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", domain.ErrInstagramTokenRefresh, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("%w: status %d: %s", domain.ErrInstagramTokenRefresh, resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil, fmt.Errorf("%w: decode: %v", domain.ErrInstagramTokenRefresh, err)
	}

	expiresAt := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return result.AccessToken, &expiresAt, nil
}

// exchangeLongLivedToken 短期トークンを長期トークン（約60日）に交換
func (c *instagramOAuthClient) exchangeLongLivedToken(ctx context.Context, shortLivedToken string) (string, *time.Time, error) {
	url := fmt.Sprintf(
		"https://graph.facebook.com/v21.0/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s",
		c.config.ClientID, c.config.ClientSecret, shortLivedToken,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("token exchange failed: status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil, err
	}

	expiresAt := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return result.AccessToken, &expiresAt, nil
}

// fetchInstagramBusinessAccount FacebookページからIGビジネスアカウントIDを取得
func (c *instagramOAuthClient) fetchInstagramBusinessAccount(ctx context.Context, accessToken string) (string, string, error) {
	url := fmt.Sprintf(
		"https://graph.facebook.com/v21.0/me/accounts?fields=id,name,instagram_business_account&access_token=%s",
		accessToken,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("failed to fetch pages: status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			ID                       string `json:"id"`
			InstagramBusinessAccount *struct {
				ID string `json:"id"`
			} `json:"instagram_business_account"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	// IGビジネスアカウントが紐づいたページを探す
	for _, page := range result.Data {
		if page.InstagramBusinessAccount != nil {
			return page.ID, page.InstagramBusinessAccount.ID, nil
		}
	}

	return "", "", fmt.Errorf("no instagram business account found on any facebook page")
}
