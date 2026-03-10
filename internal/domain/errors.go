package domain

import "errors"

// 汎用エラー
var (
	ErrNotFound        = errors.New("not found")
	ErrInternal        = errors.New("internal error")
	ErrInvalidArgument = errors.New("invalid argument")
)

// 認証・認可エラー
var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
	ErrInvalidToken       = errors.New("invalid token")
)

// OAuth関連エラー
var (
	ErrOAuthTokenExchange = errors.New("oauth token exchange failed")
	ErrOAuthUserInfo      = errors.New("failed to get oauth user info")
)

// トランザクション関連エラー
var (
	ErrTransactionAlreadyExists = errors.New("transaction already exists")
	ErrTransactionCommit        = errors.New("failed to commit transaction")
	ErrTransactionPanic         = errors.New("panic recovered in transaction")
)

// GA4 API関連エラー
var (
	ErrGA4APICall      = errors.New("ga4 api call failed")
	ErrGA4TokenRefresh = errors.New("ga4 token refresh failed")
)

// GSC API関連エラー
var (
	ErrGSCAPICall      = errors.New("gsc api call failed")
	ErrGSCTokenRefresh = errors.New("gsc token refresh failed")
)

// GBP API関連エラー
var (
	ErrGBPAPICall      = errors.New("gbp api call failed")
	ErrGBPTokenRefresh = errors.New("gbp token refresh failed")
)

// Instagram API関連エラー
var (
	ErrInstagramAPICall       = errors.New("instagram api call failed")
	ErrInstagramTokenRefresh  = errors.New("instagram token refresh failed")
	ErrInstagramTokenExchange = errors.New("instagram token exchange failed")
)

// LINE API関連エラー
var (
	ErrLineAPICall      = errors.New("line api call failed")
	ErrLineInvalidToken = errors.New("line channel access token is invalid")
)

// バリデーション関連エラー
var (
	ErrFilterRequired = errors.New("filter is required")
)
