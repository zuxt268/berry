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

type InstagramAuthHandler struct {
	instagramAuthUseCase usecase.InstagramAuthUseCase
}

func NewInstagramAuthHandler(instagramAuthUseCase usecase.InstagramAuthUseCase) *InstagramAuthHandler {
	return &InstagramAuthHandler{instagramAuthUseCase: instagramAuthUseCase}
}

// Connect Instagram OAuth連携フローを開始
func (h *InstagramAuthHandler) Connect(w http.ResponseWriter, r *http.Request) {
	url, err := h.instagramAuthUseCase.InitiateConnect(r, w)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback Instagram OAuthコールバックを処理
func (h *InstagramAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
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

	conn, err := h.instagramAuthUseCase.HandleCallback(r.Context(), r, w, user.ID, code, state)
	if err != nil {
		HandleErrorWithRedirect(w, r, err, "/settings")
		return
	}

	slog.Info("Instagram connection completed", "uid", conn.UID, "igAccountID", conn.InstagramBusinessAccountID)
	http.Redirect(w, r, config.Env.FrontendURL+"/settings?instagram=connected", http.StatusTemporaryRedirect)
}

// GetConnections 現在のユーザーのInstagram連携一覧を返す
func (h *InstagramAuthHandler) GetConnections(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleError(w, domain.ErrUnauthorized)
		return
	}

	connections, err := h.instagramAuthUseCase.GetConnections(r.Context(), user.ID)
	if err != nil {
		HandleError(w, err)
		return
	}

	responses := make([]*domain.InstagramConnectionResponse, len(connections))
	for i, c := range connections {
		responses[i] = &domain.InstagramConnectionResponse{
			UID:                        c.UID,
			InstagramBusinessAccountID: c.InstagramBusinessAccountID,
			FacebookPageID:             c.FacebookPageID,
			TokenExpiresAt:             c.TokenExpiresAt,
			ConnectedAt:                c.ConnectedAt,
			DisconnectedAt:             c.DisconnectedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"connections": responses,
	})
}

// Disconnect Instagram連携を解除
func (h *InstagramAuthHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
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

	if err := h.instagramAuthUseCase.Disconnect(r.Context(), user.ID, uid); err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "disconnected"})
}
