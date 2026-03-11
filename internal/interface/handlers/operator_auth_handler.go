package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/adapter"
	"github.com/zuxt268/berry/internal/usecase"
)

type OperatorAuthHandler struct {
	operatorAuthUseCase usecase.OperatorAuthUseCase
	sessionAdapter      adapter.SessionAdapter
}

// NewOperatorAuthHandler 新しいOperatorAuthHandlerインスタンスを作成
func NewOperatorAuthHandler(operatorAuthUseCase usecase.OperatorAuthUseCase, sessionAdapter adapter.SessionAdapter) *OperatorAuthHandler {
	return &OperatorAuthHandler{
		operatorAuthUseCase: operatorAuthUseCase,
		sessionAdapter:      sessionAdapter,
	}
}

// GoogleLogin Operator向けGoogle OAuth2ログインフローを開始
func (h *OperatorAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("operator login initiated")
	url, state, err := h.operatorAuthUseCase.InitiateLogin()
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/operator")
		return
	}
	setOAuthStateCookie(w, "operator_oauthstate", state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback Operator向けGoogle OAuth2コールバックを処理
func (h *OperatorAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	if state == "" || code == "" {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidArgument, "/operator")
		return
	}

	// state検証
	if err := verifyOAuthState(r, "operator_oauthstate", state); err != nil {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidToken, "/operator")
		return
	}
	clearCookie(w, "operator_oauthstate")

	operator, sessionToken, err := h.operatorAuthUseCase.HandleCallback(r.Context(), code, r.RemoteAddr, r.UserAgent())
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/operator")
		return
	}

	// セッショントークンをクッキーに保存
	if err := h.sessionAdapter.SaveSessionToken(r, w, sessionToken); err != nil {
		HandleErrorWithRedirect(w, r, err, "/operator")
		return
	}

	slog.Info("operator login completed", "uid", operator.UID, "email", operator.Email)
	http.Redirect(w, r, config.Env.FrontendURL+"/operator", http.StatusTemporaryRedirect)
}

// GoogleLogout Operatorログアウト処理
func (h *OperatorAuthHandler) GoogleLogout(w http.ResponseWriter, r *http.Request) {
	slog.Info("operator logout")

	// セッショントークンを取得してからusecase呼び出し
	token, _, _ := h.sessionAdapter.GetSessionToken(r)
	if err := h.operatorAuthUseCase.Logout(r.Context(), token); err != nil {
		HandleErrorWithRedirect(w, r, err, "/operator")
		return
	}

	// セッションクッキーを削除
	_ = h.sessionAdapter.DeleteSessionToken(r, w)
	clearCookie(w, "operator_oauthstate")

	http.Redirect(w, r, config.Env.FrontendURL+"/operator", http.StatusTemporaryRedirect)
}

// GetCurrentOperator 現在のオペレーターのセッション情報を返す
func (h *OperatorAuthHandler) GetCurrentOperator(w http.ResponseWriter, r *http.Request) {
	// セッショントークンを取得してからusecase呼び出し
	token, ok, err := h.sessionAdapter.GetSessionToken(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	var sessionToken string
	if ok {
		sessionToken = token
	}

	operator, authenticated, err := h.operatorAuthUseCase.GetCurrentOperator(r.Context(), sessionToken)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !authenticated || operator == nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":    operator.UID,
			"email": operator.Email,
			"name":  operator.Name,
		},
	})
}