package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/middleware"
	"github.com/zuxt268/berry/internal/usecase"
)

type LineAuthHandler struct {
	lineAuthUseCase usecase.LineAuthUseCase
}

func NewLineAuthHandler(lineAuthUseCase usecase.LineAuthUseCase) *LineAuthHandler {
	return &LineAuthHandler{lineAuthUseCase: lineAuthUseCase}
}

// Connect LINE連携を登録（POSTでトークン情報を受け取る）
func (h *LineAuthHandler) Connect(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleError(w, domain.ErrUnauthorized)
		return
	}

	var req domain.ConnectLineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	if req.ChannelID == "" || req.ChannelSecret == "" || req.ChannelAccessToken == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	conn, err := h.lineAuthUseCase.Connect(r.Context(), user.ID, &req)
	if err != nil {
		HandleError(w, err)
		return
	}

	slog.Info("LINE connection completed", "uid", conn.UID, "channelID", conn.ChannelID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&domain.LineConnectionResponse{
		UID:         conn.UID,
		ChannelID:   conn.ChannelID,
		ChannelName: conn.ChannelName,
		BotUserID:   conn.BotUserID,
		ConnectedAt: conn.ConnectedAt,
	})
}

// GetConnections 現在のユーザーのLINE連携一覧を返す
func (h *LineAuthHandler) GetConnections(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		HandleError(w, domain.ErrUnauthorized)
		return
	}

	connections, err := h.lineAuthUseCase.GetConnections(r.Context(), user.ID)
	if err != nil {
		HandleError(w, err)
		return
	}

	responses := make([]*domain.LineConnectionResponse, len(connections))
	for i, c := range connections {
		responses[i] = &domain.LineConnectionResponse{
			UID:            c.UID,
			ChannelID:      c.ChannelID,
			ChannelName:    c.ChannelName,
			BotUserID:      c.BotUserID,
			ConnectedAt:    c.ConnectedAt,
			DisconnectedAt: c.DisconnectedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"connections": responses,
	})
}

// Disconnect LINE連携を解除
func (h *LineAuthHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
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

	if err := h.lineAuthUseCase.Disconnect(r.Context(), user.ID, uid); err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "disconnected"})
}
