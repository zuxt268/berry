package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase"
)

type OperatorHandler struct {
	operatorUsecase usecase.OperatorUsecase
}

func NewOperatorHandler(operatorUsecase usecase.OperatorUsecase) *OperatorHandler {
	return &OperatorHandler{operatorUsecase: operatorUsecase}
}

// GetByUID GET /api/operators/{uid}
func (h *OperatorHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	operator, err := h.operatorUsecase.GetByUID(r.Context(), uid)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, operator)
}

// Gets GET /api/operators
func (h *OperatorHandler) Gets(w http.ResponseWriter, r *http.Request) {
	var input domain.GetOperatorsRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		input = domain.GetOperatorsRequest{}
	}

	resp, err := h.operatorUsecase.Gets(r.Context(), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// Create POST /api/operators
func (h *OperatorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateOperatorRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	if err := Validate(input); err != nil {
		HandleError(w, err)
		return
	}

	resp, err := h.operatorUsecase.Create(r.Context(), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// Update PUT /api/operators/{uid}
func (h *OperatorHandler) Update(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	var input domain.UpdateOperatorRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}
	input.UID = uid

	resp, err := h.operatorUsecase.Update(r.Context(), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// Delete DELETE /api/operators/{uid}
func (h *OperatorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, domain.ErrInvalidArgument)
		return
	}

	if err := h.operatorUsecase.Delete(r.Context(), uid); err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
