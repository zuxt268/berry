package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
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

	url, err := h.gbpAuthUseCase.InitiateConnect(r, w, locationID, accountID)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

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

	conn, err := h.gbpAuthUseCase.HandleCallback(r.Context(), r, w, user.ID, code, state)
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

	responses := make([]*domain.GBPConnectionResponse, len(connections))
	for i, c := range connections {
		responses[i] = &domain.GBPConnectionResponse{
			UID:            c.UID,
			LocationID:     c.LocationID,
			AccountID:      c.AccountID,
			ConnectedAt:    c.ConnectedAt,
			DisconnectedAt: c.DisconnectedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"connections": responses,
	})
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "disconnected"})
}
