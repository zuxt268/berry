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

type GBPAuthHandler struct {
	gbpAuthUseCase usecase.GBPAuthUseCase
}

func NewGBPAuthHandler(gbpAuthUseCase usecase.GBPAuthUseCase) *GBPAuthHandler {
	return &GBPAuthHandler{gbpAuthUseCase: gbpAuthUseCase}
}

// Connect GBP OAuth連携フローを開始
func (h *GBPAuthHandler) Connect(w http.ResponseWriter, r *http.Request) {
	locationID := r.URL.Query().Get("location_id")
	accountID := r.URL.Query().Get("account_id")
	if locationID == "" || accountID == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	url, state, err := h.gbpAuthUseCase.InitiateConnect(locationID, accountID)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	// stateとlocation_id、account_idをクッキーに保存
	setOAuthStateCookie(w, "gbp_oauthstate", state)
	setValueCookie(w, "gbp_location_id", locationID)
	setValueCookie(w, "gbp_account_id", accountID)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback GBP OAuthコールバックを処理
func (h *GBPAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
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
	if err := verifyOAuthState(r, "gbp_oauthstate", state); err != nil {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidToken, "/settings")
		return
	}

	// location_id取得
	locationID, err := getOAuthStateCookie(r, "gbp_location_id")
	if err != nil || locationID == "" {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidArgument, "/settings")
		return
	}

	// account_id取得
	accountID, err := getOAuthStateCookie(r, "gbp_account_id")
	if err != nil || accountID == "" {
		HandleErrorWithRedirect(w, r, domain.ErrInvalidArgument, "/settings")
		return
	}

	// クッキー削除
	clearCookie(w, "gbp_oauthstate")
	clearCookie(w, "gbp_location_id")
	clearCookie(w, "gbp_account_id")

	conn, err := h.gbpAuthUseCase.HandleCallback(r.Context(), user.ID, code, locationID, accountID)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	slog.Info("GBP connection completed", "uid", conn.UID, "locationID", conn.LocationID)
	http.Redirect(w, r, config.Env.FrontendURL+"/settings?gbp=connected", http.StatusTemporaryRedirect)
}

// GetConnections 現在のユーザーのGBP連携一覧を返す
func (h *GBPAuthHandler) GetConnections(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleError(w, domain.ErrUnauthorized)
		return
	}

	connections, err := h.gbpAuthUseCase.GetConnections(r.Context(), user.ID)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp := make([]*responses.GBPConnectionResponse, len(connections))
	for i, c := range connections {
		resp[i] = responses.ToGBPConnectionResponse(c)
	}

	respondJSON(w, http.StatusOK, map[string]any{"connections": resp})
}

// Disconnect GBP連携を解除
func (h *GBPAuthHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
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

	if err := h.gbpAuthUseCase.Disconnect(r.Context(), user.ID, uid); err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}