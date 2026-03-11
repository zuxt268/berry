package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/dto/responses"
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

	var req usecase.ConnectLineInput
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

	respondJSON(w, http.StatusCreated, responses.ToLineConnectionResponse(conn))
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

	resp := make([]*responses.LineConnectionResponse, len(connections))
	for i, c := range connections {
		resp[i] = responses.ToLineConnectionResponse(c)
	}

	respondJSON(w, http.StatusOK, map[string]any{"connections": resp})
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

	respondJSON(w, http.StatusOK, map[string]string{"status": "disconnected"})
}