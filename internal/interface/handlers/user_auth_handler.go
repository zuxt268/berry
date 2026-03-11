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

type UserAuthHandler struct {
	userAuthUseCase usecase.UserAuthUseCase
	sessionAdapter  adapter.SessionAdapter
}

// NewUserAuthHandler 新しいAuthHandlerインスタンスを作成
func NewUserAuthHandler(userAuthUseCase usecase.UserAuthUseCase, sessionAdapter adapter.SessionAdapter) *UserAuthHandler {
	return &UserAuthHandler{
		userAuthUseCase: userAuthUseCase,
		sessionAdapter:  sessionAdapter,
	}
}

// GoogleLogin Google OAuth2ログインフローを開始
func (h *UserAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("customer login initiated")
	url, state, err := h.userAuthUseCase.InitiateLogin()
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/")
		return
	}
	setOAuthStateCookie(w, "oauthstate", state)
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

	// state検証
	if err := verifyOAuthState(r, "oauthstate", state); err != nil {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidToken, "/")
		return
	}
	clearCookie(w, "oauthstate")

	user, sessionToken, err := h.userAuthUseCase.HandleCallback(r.Context(), code, r.RemoteAddr, r.UserAgent())
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/")
		return
	}

	// セッショントークンをクッキーに保存
	if err := h.sessionAdapter.SaveSessionToken(r, w, sessionToken); err != nil {
		HandleErrorWithRedirect(w, r, err, "/")
		return
	}

	slog.Info("customer login completed", "uid", user.UID, "email", user.Email)
	http.Redirect(w, r, config.Env.FrontendURL+"/", http.StatusTemporaryRedirect)
}

// GoogleLogout ログアウト処理
func (h *UserAuthHandler) GoogleLogout(w http.ResponseWriter, r *http.Request) {
	slog.Info("customer logout")

	// セッショントークンを取得してからusecase呼び出し
	token, _, _ := h.sessionAdapter.GetSessionToken(r)
	if err := h.userAuthUseCase.Logout(r.Context(), token); err != nil {
		HandleErrorWithRedirect(w, r, err, "/")
		return
	}

	// セッションクッキーを削除
	_ = h.sessionAdapter.DeleteSessionToken(r, w)
	clearCookie(w, "oauthstate")

	http.Redirect(w, r, config.Env.FrontendURL+"/", http.StatusTemporaryRedirect)
}

// GetCurrentUser 現在のユーザーのセッション情報を返す
func (h *UserAuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
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

	user, authenticated, err := h.userAuthUseCase.GetCurrentUser(r.Context(), sessionToken)
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