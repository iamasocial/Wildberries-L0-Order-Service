package http

import (
	"L0/internal/service"
	"L0/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *service.Service
	logger  logger.Logger
}

func NewHandler(service *service.Service, logger logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderUID := chi.URLParam(r, "orderUID")
	h.logger.Debug("received GET /order request", "order_uid", orderUID)

	order, err := h.service.Order.GetOrderByUID(ctx, orderUID)
	if err != nil {
		h.logger.Error("failed to get order", "order_uid", orderUID, "error", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// h.logger.Info("order re", "order_uid", orderUID)
	h.logger.Info("order retrieved successfully", "order_uid", orderUID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
