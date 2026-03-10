package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase"
)

type UserAuthHandler struct {
	userAuthUseCase usecase.UserAuthUseCase
}

// NewUserAuthHandler 新しいAuthHandlerインスタンスを作成
func NewUserAuthHandler(userAuthUseCase usecase.UserAuthUseCase) *UserAuthHandler {
	return &UserAuthHandler{
		userAuthUseCase: userAuthUseCase,
	}
}

// GoogleLogin Google OAuth2ログインフローを開始
func (h *UserAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("customer login initiated")
	url, err := h.userAuthUseCase.InitiateLogin(w)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/")
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback Google OAuth2コールバックを処理
func (h *UserAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	if state == "" || code == "" {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidArgument, "/")
		return
	}

	user, err := h.userAuthUseCase.HandleCallback(r.Context(), r, w, code, state)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/")
		return
	}

	slog.Info("customer login completed", "uid", user.UID, "email", user.Email)
	http.Redirect(w, r, config.Env.FrontendURL+"/", http.StatusTemporaryRedirect)
}

// GoogleLogout ログアウト処理
func (h *UserAuthHandler) GoogleLogout(w http.ResponseWriter, r *http.Request) {
	slog.Info("customer logout")
	if err := h.userAuthUseCase.Logout(r.Context(), r, w); err != nil {
		HandleErrorWithRedirect(w, r, err, "/")
		return
	}

	cookie := &http.Cookie{
		Name:     "oauthstate",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.Env.AppEnv == "prod",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, config.Env.FrontendURL+"/", http.StatusTemporaryRedirect)
}

// GetCurrentUser 現在のユーザーのセッション情報を返す
func (h *UserAuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, authenticated, err := h.userAuthUseCase.GetCurrentUser(r.Context(), r)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !authenticated || user == nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":    user.UID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
