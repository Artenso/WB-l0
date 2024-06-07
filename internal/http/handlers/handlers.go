package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Artenso/wb-l0/internal/service"
)

// IHandler working with handler
type IHandler interface {
	GetOrder(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service service.IService
}

// New creates new handler
func New(service service.IService) IHandler {
	return &handler{
		service: service,
	}
}

// GetOrder gets order by order_uid
func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Path[len("/order/"):]
	if len(orderUID) == 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
