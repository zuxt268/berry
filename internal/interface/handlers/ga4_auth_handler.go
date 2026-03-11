package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/dto/responses"
	"github.com/zuxt268/berry/internal/interface/middleware"
	"github.com/zuxt268/berry/internal/usecase"
)

type GA4AuthHandler struct {
	ga4AuthUseCase usecase.GA4AuthUseCase
}

func NewGA4AuthHandler(ga4AuthUseCase usecase.GA4AuthUseCase) *GA4AuthHandler {
	return &GA4AuthHandler{ga4AuthUseCase: ga4AuthUseCase}
}

// Connect GA4 OAuth連携フローを開始
func (h *GA4AuthHandler) Connect(w http.ResponseWriter, r *http.Request) {
	propertyID := r.URL.Query().Get("property_id")
	if propertyID == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	url, state, err := h.ga4AuthUseCase.InitiateConnect(propertyID)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	// stateとproperty_idをクッキーに保存
	setOAuthStateCookie(w, "ga4_oauthstate", state)
	setValueCookie(w, "ga4_property_id", propertyID)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback GA4 OAuthコールバックを処理
func (h *GA4AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	if state == "" || code == "" {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidArgument, "/settings")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleErrorWithRedirect(w, r, domain.ErrUnauthorized, "/settings")
		return
	}

	// state検証
	if err := verifyOAuthState(r, "ga4_oauthstate", state); err != nil {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidToken, "/settings")
		return
	}

	// property_id取得
	propertyID, err := getOAuthStateCookie(r, "ga4_property_id")
	if err != nil || propertyID == "" {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidArgument, "/settings")
		return
	}

	// クッキー削除
	clearCookie(w, "ga4_oauthstate")
	clearCookie(w, "ga4_property_id")

	conn, err := h.ga4AuthUseCase.HandleCallback(r.Context(), user.ID, code, propertyID)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	slog.Info("GA4 connection completed", "uid", conn.UID, "propertyID", conn.GooglePropertyID)
	http.Redirect(w, r, config.Env.FrontendURL+"/settings?ga4=connected", http.StatusTemporaryRedirect)
}

// GetConnections 現在のユーザーのGA4連携一覧を返す
func (h *GA4AuthHandler) GetConnections(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleError(w, domain.ErrUnauthorized)
		return
	}

	connections, err := h.ga4AuthUseCase.GetConnections(r.Context(), user.ID)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp := make([]*responses.GA4ConnectionResponse, len(connections))
	for i, c := range connections {
		resp[i] = responses.ToGA4ConnectionResponse(c)
	}

	respondJSON(w, http.StatusOK, map[string]any{"connections": resp})
}

// Disconnect GA4連携を解除
func (h *GA4AuthHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleError(w, domain.ErrUnauthorized)
		return
	}

	if err := h.ga4AuthUseCase.Disconnect(r.Context(), user.ID, uid); err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}