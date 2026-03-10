package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase"
)

type OperatorAuthHandler struct {
	operatorAuthUseCase usecase.OperatorAuthUseCase
}

// NewOperatorAuthHandler 新しいOperatorAuthHandlerインスタンスを作成
func NewOperatorAuthHandler(operatorAuthUseCase usecase.OperatorAuthUseCase) *OperatorAuthHandler {
	return &OperatorAuthHandler{
		operatorAuthUseCase: operatorAuthUseCase,
	}
}

// GoogleLogin Operator向けGoogle OAuth2ログインフローを開始
func (h *OperatorAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("operator login initiated")
	url, err := h.operatorAuthUseCase.InitiateLogin(w)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/operator")
		return
	}
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

	operator, err := h.operatorAuthUseCase.HandleCallback(r.Context(), r, w, code, state)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/operator")
		return
	}

	slog.Info("operator login completed", "uid", operator.UID, "email", operator.Email)
	http.Redirect(w, r, config.Env.FrontendURL+"/operator", http.StatusTemporaryRedirect)
}

// GoogleLogout Operatorログアウト処理
func (h *OperatorAuthHandler) GoogleLogout(w http.ResponseWriter, r *http.Request) {
	slog.Info("operator logout")
	if err := h.operatorAuthUseCase.Logout(r.Context(), r, w); err != nil {
		HandleErrorWithRedirect(w, r, err, "/operator")
		return
	}

	cookie := &http.Cookie{
		Name:     "operator_oauthstate",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.Env.AppEnv == "prod",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, config.Env.FrontendURL+"/operator", http.StatusTemporaryRedirect)
}

// GetCurrentOperator 現在のオペレーターのセッション情報を返す
func (h *OperatorAuthHandler) GetCurrentOperator(w http.ResponseWriter, r *http.Request) {
	operator, authenticated, err := h.operatorAuthUseCase.GetCurrentOperator(r.Context(), r)
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
