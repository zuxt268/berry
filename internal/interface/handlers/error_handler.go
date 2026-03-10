package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
)

// ErrorResponse エラーレスポンスの構造体
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// getErrorInfo エラーから情報を取得
func getErrorInfo(err error) (statusCode int, errorType string, message string) {
	if err == nil {
		return http.StatusInternalServerError, "unknown_error", "予期しないエラーが発生しました"
	}

	switch {
	// 認証・認可エラー
	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		errorType = "unauthorized"
		message = "認証が必要です"
	case errors.Is(err, domain.ErrForbidden):
		statusCode = http.StatusForbidden
		errorType = "forbidden"
		message = "アクセス権限がありません"
	case errors.Is(err, domain.ErrInvalidCredentials):
		statusCode = http.StatusUnauthorized
		errorType = "invalid_credentials"
		message = "認証情報が無効です"
	case errors.Is(err, domain.ErrSessionExpired):
		statusCode = http.StatusUnauthorized
		errorType = "session_expired"
		message = "セッションが期限切れです"
	case errors.Is(err, domain.ErrInvalidToken):
		statusCode = http.StatusUnauthorized
		errorType = "invalid_token"
		message = "トークンが無効です"

	// OAuth関連エラー
	case errors.Is(err, domain.ErrOAuthTokenExchange):
		statusCode = http.StatusBadGateway
		errorType = "oauth_token_exchange"
		message = "トークン交換に失敗しました"
	case errors.Is(err, domain.ErrOAuthUserInfo):
		statusCode = http.StatusBadGateway
		errorType = "oauth_user_info"
		message = "ユーザー情報の取得に失敗しました"

	// バリデーション関連エラー
	case errors.Is(err, domain.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		errorType = "invalid_argument"
		message = "引数が不正です"
	case errors.Is(err, domain.ErrFilterRequired):
		statusCode = http.StatusBadRequest
		errorType = "filter_required"
		message = "フィルターが必要です"

	// 汎用エラー
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		errorType = "not_found"
		message = "リソースが見つかりません"
	case errors.Is(err, domain.ErrInternal):
		statusCode = http.StatusInternalServerError
		errorType = "internal_error"
		message = "内部エラーが発生しました"

	// トランザクション関連エラー
	case errors.Is(err, domain.ErrTransactionAlreadyExists):
		statusCode = http.StatusConflict
		errorType = "transaction_already_exists"
		message = "トランザクションが既に存在します"
	case errors.Is(err, domain.ErrTransactionCommit):
		statusCode = http.StatusInternalServerError
		errorType = "transaction_commit"
		message = "トランザクションのコミットに失敗しました"
	case errors.Is(err, domain.ErrTransactionPanic):
		statusCode = http.StatusInternalServerError
		errorType = "transaction_panic"
		message = "トランザクション中にパニックが発生しました"

	// デフォルト
	default:
		statusCode = http.StatusInternalServerError
		errorType = "unknown_error"
		message = "予期しないエラーが発生しました"
	}

	return statusCode, errorType, message
}

// HandleError エラーを適切なHTTPレスポンス(JSON)に変換
func HandleError(w http.ResponseWriter, err error) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		HandleValidationError(w, ve)
		return
	}

	statusCode, errorType, message := getErrorInfo(err)

	errMsg := "nil"
	if err != nil {
		errMsg = err.Error()
	}
	slog.Error("エラー発生",
		"error", errMsg,
		"type", errorType,
		"status", statusCode,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}

// HandleErrorWithRedirect エラーをリダイレクト先ページへのエラーパラメータ付きリダイレクトに変換
func HandleErrorWithRedirect(w http.ResponseWriter, r *http.Request, err error, redirectPath string) {
	statusCode, errorType, _ := getErrorInfo(err)

	errMsg := "nil"
	if err != nil {
		errMsg = err.Error()
	}
	slog.Error("エラー発生",
		"error", errMsg,
		"type", errorType,
		"status", statusCode,
	)

	redirectURL := config.Env.FrontendURL + redirectPath + "?error=" + errorType
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
