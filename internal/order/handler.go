package order

import (
	"codealong/internal/platform/httpx"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", h.placeOrder)
	mux.HandleFunc("GET /orders/{id}", h.getOrder)
}

func (h *Handler) placeOrder(w http.ResponseWriter, r *http.Request) {
	var req PlaceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	order, err := h.svc.PlaceOrder(r.Context(), req)
	switch {
	case err == nil:
		httpx.JSON(w, http.StatusCreated, order)
	case errors.Is(err, ErrCustomerRequired), errors.Is(err, ErrNoItems):
		httpx.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrCustomerNotFound), errors.Is(err, ErrProductNotFound):
		httpx.Error(w, http.StatusUnprocessableEntity, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	safeID, err := uuid.Parse(id)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	order, err := h.svc.Get(ctx, safeID.String())
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.Error(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, order)
}
