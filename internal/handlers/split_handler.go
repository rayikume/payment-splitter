package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/rayikume/payment-splitter/internal/models"
	"github.com/rayikume/payment-splitter/internal/responses"
	"github.com/rayikume/payment-splitter/internal/services"
)

type SplitHandler struct {
	svc *services.SplitService
}

func NewSplitHandler(svc *services.SplitService) *SplitHandler {
	return &SplitHandler{svc: svc}
}

func (p *SplitHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/splits", p.Create)

}

func (p *SplitHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSplitRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "INVALID_JSON", "request body is not valid JSON")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		responses.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "created_by is required")
		return
	}
	if req.Currency == "" {
		responses.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "currency is required")
		return
	}

	split, err := p.svc.Create(req)
	if err != nil {
		mapServiceError(w, err)
		return
	}

	responses.Success(w, http.StatusCreated, split)
}

func (p *SplitHandler) GetByID(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	if id == "" {
		responses.Error(w, http.StatusBadRequest, "INVALID_PARAM", "split id is required")
		return
	}
	split, err := p.svc.GetByID(id)
	if err != nil {
		mapServiceError(w, err)
		return
	}

	responses.Success(w, http.StatusOK, split)
}

func (p *SplitHandler) Settle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req models.SettleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, "INVALID_JSON", "request body is not valid JSON")
		return
	}
	if strings.TrimSpace(req.ParticipantID) == "" {
		responses.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "participant_id is required")
		return
	}

	split, err := p.svc.Settle(id, req.ParticipantID)
	if err != nil {
		mapServiceError(w, err)
		return
	}

	responses.Success(w, http.StatusOK, split)
}

func (p *SplitHandler) Delete(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	if err := p.svc.Delete(id); err != nil {
		mapServiceError(w, err)
		return
	}

	responses.Success(w, http.StatusNoContent, nil)
}

func mapServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrSplitNotFound):
		responses.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, services.ErrParticipantNotFound):
		responses.Error(w, http.StatusNotFound, "PARTICIPANT_NOT_FOUND", err.Error())
	case errors.Is(err, services.ErrInvalidStrategy),
		errors.Is(err, services.ErrAmountMismatch),
		errors.Is(err, services.ErrPercentageMismatch),
		errors.Is(err, services.ErrLessThanMinimumParticipants),
		errors.Is(err, services.ErrInvalidAmount):
		responses.Error(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
	default:
		responses.Error(w, http.StatusInternalServerError, "INTERNAL", "an unexpected error occurred")
	}
}
