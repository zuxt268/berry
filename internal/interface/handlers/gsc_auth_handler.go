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

type GSCAuthHandler struct {
	gscAuthUseCase usecase.GSCAuthUseCase
}

func NewGSCAuthHandler(gscAuthUseCase usecase.GSCAuthUseCase) *GSCAuthHandler {
	return &GSCAuthHandler{gscAuthUseCase: gscAuthUseCase}
}

// Connect GSC OAuth連携フローを開始
func (h *GSCAuthHandler) Connect(w http.ResponseWriter, r *http.Request) {
	siteURL := r.URL.Query().Get("site_url")
	if siteURL == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	url, err := h.gscAuthUseCase.InitiateConnect(r, w, siteURL)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback GSC OAuthコールバックを処理
func (h *GSCAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
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

	conn, err := h.gscAuthUseCase.HandleCallback(r.Context(), r, w, user.ID, code, state)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	slog.Info("GSC connection completed", "uid", conn.UID, "siteURL", conn.SiteURL)
	http.Redirect(w, r, config.Env.FrontendURL+"/settings?gsc=connected", http.StatusTemporaryRedirect)
}

// GetConnections 現在のユーザーのGSC連携一覧を返す
func (h *GSCAuthHandler) GetConnections(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleError(w, domain.ErrUnauthorized)
		return
	}

	connections, err := h.gscAuthUseCase.GetConnections(r.Context(), user.ID)
	if err != nil {
		HandleError(w, err)
		return
	}

	responses := make([]*domain.GSCConnectionResponse, len(connections))
	for i, c := range connections {
		responses[i] = &domain.GSCConnectionResponse{
			UID:            c.UID,
			SiteURL:        c.SiteURL,
			ConnectedAt:    c.ConnectedAt,
			DisconnectedAt: c.DisconnectedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"connections": responses,
	})
}

// Disconnect GSC連携を解除
func (h *GSCAuthHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
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

	if err := h.gscAuthUseCase.Disconnect(r.Context(), user.ID, uid); err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "disconnected"})
}
