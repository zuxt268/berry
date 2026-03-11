package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/dto/responses"
	"github.com/zuxt268/berry/internal/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// GetByUID GET /api/users/{uid}
func (h *UserHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	user, err := h.userUsecase.GetByUID(r.Context(), uid)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, responses.ToUserResponse(user))
}

// Gets GET /api/users
func (h *UserHandler) Gets(w http.ResponseWriter, r *http.Request) {
	var input usecase.GetUsersInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// クエリパラメータのみの場合はデコードエラーを無視
		input = usecase.GetUsersInput{}
	}

	users, total, err := h.userUsecase.Gets(r.Context(), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, responses.ToUsersResponse(users, total))
}

// Create POST /api/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	if err := Validate(input); err != nil {
		HandleError(w, err)
		return
	}

	user, err := h.userUsecase.Create(r.Context(), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, responses.ToUserResponse(user))
}

// Update PUT /api/users/{uid}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	var input usecase.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}
	input.UID = uid

	user, err := h.userUsecase.Update(r.Context(), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, responses.ToUserResponse(user))
}

// Delete DELETE /api/users/{uid}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	if err := h.userUsecase.Delete(r.Context(), uid); err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}