package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthUser OAuth認証で取得したユーザー情報
type OAuthUser struct {
	Sub     string
	Email   string
	Name    string
	Picture string
}

// OAuthAdapter OAuth操作のインターフェース
//
//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package mock
type OAuthAdapter interface {
	// GetAuthURL stateパラメータ付きのOAuth認証URLを返す
	GetAuthURL(state string) string
	// ExchangeCode 認証コードをユーザー情報と交換
	ExchangeCode(ctx context.Context, code string) (*OAuthUser, error)
}

type oauthClient struct {
	config *oauth2.Config
}

func NewOAuthClient(redirectURL string) OAuthAdapter {
	oauthConfig := &oauth2.Config{
		ClientID:     config.Env.GoogleClientID,
		ClientSecret: config.Env.GoogleClientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &oauthClient{config: oauthConfig}
}

// GetAuthURL stateパラメータ付きのOAuth認証URLを返す
func (c *oauthClient) GetAuthURL(state string) string {
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode 認証コードをユーザー情報と交換
func (c *oauthClient) ExchangeCode(ctx context.Context, code string) (*OAuthUser, error) {
	// コードをトークンに交換
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrOAuthTokenExchange, err)
	}

	// ユーザー情報を取得
	userInfo, err := c.getUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrOAuthUserInfo, err)
	}

	return &OAuthUser{
		Sub:     getStringFromMap(userInfo, "id"),
		Email:   getStringFromMap(userInfo, "email"),
		Name:    getStringFromMap(userInfo, "name"),
		Picture: getStringFromMap(userInfo, "picture"),
	}, nil
}

func (c *oauthClient) getUserInfo(accessToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", domain.ErrOAuthUserInfo, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
