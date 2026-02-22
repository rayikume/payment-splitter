package handlers

import (
	"net/http"

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

}
